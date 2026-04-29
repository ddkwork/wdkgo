package main

type SerialBaudRate = uint32

const (
	Baud110    SerialBaudRate = 110
	Baud300    SerialBaudRate = 300
	Baud600    SerialBaudRate = 600
	Baud1200   SerialBaudRate = 1200
	Baud2400   SerialBaudRate = 2400
	Baud4800   SerialBaudRate = 4800
	Baud9600   SerialBaudRate = 9600
	Baud14400  SerialBaudRate = 14400
	Baud19200  SerialBaudRate = 19200
	Baud38400  SerialBaudRate = 38400
	Baud57600  SerialBaudRate = 57600
	Baud115200 SerialBaudRate = 115200
)

type ParityType = int

const (
	ParityNone  ParityType = 0
	ParityOdd   ParityType = 1
	ParityEven  ParityType = 2
	ParityMark  ParityType = 3
	ParitySpace ParityType = 4
)

type StopBits = int

const (
	StopBits1   StopBits = 0
	StopBits1_5 StopBits = 1
	StopBits2   StopBits = 2
)

type SerialConfig struct {
	PortNumber uint32
	BaudRate   SerialBaudRate
	DataBits   uint32
	Parity     ParityType
	StopBits   StopBits
	UseIrq     bool
	IrqNumber  uint32
}

type KdSerialState struct {
	config       SerialConfig
	connected    bool
	lock         *Spinlock
	txBuffer     []byte
	rxBuffer     []byte
	txHead       uint32
	txTail       uint32
	rxHead       uint32
	rxTail       uint32
	bufferSize   uint32
	portBase     uint16
	lineControl  uint8
	modemControl uint8
	intEnable    uint8
	divisorLatch bool
}

func NewKdSerialState() *KdSerialState {
	bufSize := uint32(4096)
	s := &KdSerialState{
		config: SerialConfig{
			PortNumber: 2,
			BaudRate:   Baud115200,
			DataBits:   8,
			Parity:     ParityNone,
			StopBits:   StopBits1,
			UseIrq:     false,
		},
		lock:        NewSpinlock(),
		txBuffer:    make([]byte, bufSize),
		rxBuffer:    make([]byte, bufSize),
		bufferSize:  bufSize,
		lineControl: 0x03,
		portBase:    0x2F8,
	}
	return s
}

func (s *KdSerialState) Initialize(config *SerialConfig) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if config != nil {
		s.config = *config
	}

	s.portBase = s.comPortBase(s.config.PortNumber)
	if s.portBase == 0 {
		return fmtError("invalid COM port number: %d", s.config.PortNumber)
	}

	if err := s.configureHardware(); err != nil {
		return err
	}

	s.connected = true
	LogInfo("KD serial initialized on COM%d @0x%X, %d baud",
		s.config.PortNumber, s.portBase, s.config.BaudRate)
	return nil
}

func (s *KdSerialState) Uninitialize() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.connected {
		return nil
	}

	s.disableUart()
	s.connected = false
	s.txHead = 0
	s.txTail = 0
	s.rxHead = 0
	s.rxTail = 0

	LogInfo("KD serial uninitialized")
	return nil
}

func (s *KdSerialState) IsConnected() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.connected
}

func (s *KdSerialState) comPortBase(portNum uint32) uint16 {
	switch portNum {
	case 1:
		return 0x3F8
	case 2:
		return 0x2F8
	case 3:
		return 0x3E8
	case 4:
		return 0x2E8
	default:
		return 0
	}
}

func (s *KdSerialState) configureHardware() error {

	s.writePort(UART_IER, 0)

	s.setDlab(true)

	divisor := s.calculateDivisor()
	s.writePort(UART_DLL, uint8(divisor&0xFF))
	s.writePort(UART_DLM, uint8((divisor>>8)&0xFF))

	s.setDlab(false)

	lcr := uint8(0x03)
	switch s.config.Parity {
	case ParityOdd:
		lcr |= 0x08
	case ParityEven:
		lcr |= 0x18
	case ParityMark:
		lcr |= 0x28
	case ParitySpace:
		lcr |= 0x38
	}

	switch s.config.DataBits {
	case 5:
		lcr &= ^uint8(0x03)
	case 6:
		lcr = (lcr & ^uint8(0x03)) | 0x01
	case 7:
		lcr = (lcr & ^uint8(0x03)) | 0x02
	default:
		lcr = (lcr & ^uint8(0x03)) | 0x03
	}

	switch s.config.StopBits {
	case StopBits2:
		lcr |= 0x04
	}

	s.lineControl = lcr
	s.writePort(UART_LCR, lcr)

	s.writePort(UART_FCR, 0xC7)

	mcr := uint8(0x0B)
	s.modemControl = mcr
	s.writePort(UART_MCR, mcr)

	s.writePort(UART_IER, 0x00)

	s.readPort(UART_RBR)
	s.readPort(UART_IIR)
	s.readPort(UART_LSR)
	s.readPort(UART_MSR)

	LogDebug("SERIAL: Hardware configured, divisor=%d, LCR=0x%02X", divisor, lcr)
	return nil
}

func (s *KdSerialState) calculateDivisor() uint16 {
	const baseClock uint32 = 115200
	rate := uint32(s.config.BaudRate)
	if rate == 0 {
		rate = 9600
	}
	return uint16(baseClock / rate)
}

func (s *KdSerialState) setDlab(enable bool) {
	if enable {
		s.lineControl |= 0x80
	} else {
		s.lineControl &= ^uint8(0x80)
	}
	s.writePort(UART_LCR, s.lineControl)
	s.divisorLatch = enable
}

func (s *KdSerialState) disableUart() {
	s.writePort(UART_IER, 0x00)
	s.writePort(UART_FCR, 0x00)
	s.writePort(UART_MCR, 0x00)
}

func (s *KdSerialState) writePort(reg uint8, value uint8) {
	outb(s.portBase+uint16(reg), value)
}

func (s *KdSerialState) readPort(reg uint8) uint8 {
	return inb(s.portBase + uint16(reg))
}

func (s *KdSerialState) Send(data []byte) (int, error) {
	if !s.IsConnected() {
		return 0, fmtError("serial not connected")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	bytesWritten := 0
	for _, b := range data {
		nextTail := (s.txTail + 1) % s.bufferSize
		if nextTail == s.txHead {
			break
		}
		s.txBuffer[s.txTail] = b
		s.txTail = nextTail
		bytesWritten++
	}

	if bytesWritten > 0 {
		s.flushTxBufferHardware()
	}

	return bytesWritten, nil
}

func (s *KdSerialState) flushTxBufferHardware() {
	for s.txHead != s.txTail {
		if !s.isTransmitEmpty() {
			continue
		}
		b := s.txBuffer[s.txHead]
		s.writePort(UART_THR, b)
		s.txHead = (s.txHead + 1) % s.bufferSize
	}
}

func (s *KdSerialState) isTransmitEmpty() bool {
	lsr := s.readPort(UART_LSR)
	return lsr&0x20 != 0
}

func (s *KdSerialState) isReceiveReady() bool {
	lsr := s.readPort(UART_LSR)
	return lsr&0x01 != 0
}

func (s *KdSerialState) Receive(buffer []byte) (int, error) {
	if !s.IsConnected() {
		return 0, fmtError("serial not connected")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	available := s.bytesInRxBuffer()
	if available == 0 {
		s.pollRxHardware()
		available = s.bytesInRxBuffer()
	}

	bytesRead := 0
	maxRead := min(len(buffer), int(available))

	for i := range maxRead {
		if s.rxHead == s.rxTail {
			break
		}
		buffer[i] = s.rxBuffer[s.rxHead]
		s.rxHead = (s.rxHead + 1) % s.bufferSize
		bytesRead++
	}

	return bytesRead, nil
}

func (s *KdSerialState) pollRxHardware() {
	for s.isReceiveReady() {
		nextTail := (s.rxTail + 1) % s.bufferSize
		if nextTail == s.rxHead {
			break
		}
		b := s.readPort(UART_RBR)
		s.rxBuffer[s.rxTail] = b
		s.rxTail = nextTail
	}
}

func (s *KdSerialState) bytesInRxBuffer() uint32 {
	if s.rxTail >= s.rxHead {
		return s.rxTail - s.rxHead
	}
	return s.bufferSize - s.rxHead + s.rxTail
}

func (s *KdSerialState) txBytesQueued() uint32 {
	if s.txHead > s.txTail {
		return s.bufferSize - s.txHead + s.txTail
	}
	return s.txTail - s.txHead
}

func (s *KdSerialState) SendString(str string) (int, error) {
	return s.Send([]byte(str))
}

func (s *KdSerialState) ReceiveString(maxLen int) (string, error) {
	buf := make([]byte, maxLen)
	n, err := s.Receive(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func (s *KdSerialState) BytesAvailable() uint32 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pollRxHardware()
	return s.bytesInRxBuffer()
}

func (s *KdSerialState) TxSpaceAvailable() uint32 {
	s.lock.Lock()
	defer s.lock.Unlock()
	queued := s.txBytesQueued()
	if queued >= s.bufferSize-1 {
		return 0
	}
	return s.bufferSize - 1 - queued
}

func (s *KdSerialState) SendByte(b byte) error {
	data := []byte{b}
	_, err := s.Send(data)
	return err
}

func (s *KdSerialState) ReceiveByte() (byte, error) {
	buf := make([]byte, 1)
	n, err := s.Receive(buf)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmtError("no data available")
	}
	return buf[0], nil
}

const (
	UART_RBR uint8 = 0
	UART_THR uint8 = 0
	UART_IER uint8 = 1
	UART_DLL uint8 = 0
	UART_DLM uint8 = 1
	UART_IIR uint8 = 2
	UART_FCR uint8 = 2
	UART_LCR uint8 = 3
	UART_MCR uint8 = 4
	UART_LSR uint8 = 5
	UART_MSR uint8 = 6
	UART_SCR uint8 = 7
)
