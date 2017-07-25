//go:generate go-bindata -prefix ir-data ./ir-data

package main

import . "github.com/deadsy/libusb"
import "fmt"

// USB IR remocon ADVANCE
const (
	VENDOR_ID  = 0x22ea
	PRODUCT_ID = 0x003a
	BTO_EP_IN  = 0x84
	BTO_EP_OUT = 0x04
)

const (
	IR_FREQ_DEFAULT               = 38000
	IR_FREQ_MIN                   = 25000 // 赤外線周波数設定最小値 25KHz
	IR_FREQ_MAX                   = 50000 // 赤外線周波数設定最大値 50KHz
	IR_SEND_DATA_USB_SEND_MAX_LEN = 14    // USB送信１回で送信する最大ビット数
	IR_SEND_DATA_MAX_LEN          = 300   // 赤外線送信データ設定最大長[byte]
	IR_SEND_DATA_MAX_LEN_K        = IR_SEND_DATA_MAX_LEN * 8
)

type IR struct {
	CTX Context
	DH  Device_Handle
}

func CreateIR() (*IR, bool) {
	i := IR{CTX: nil, DH: nil}

	err := Init(&i.CTX)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	return &i, true
}

func (i *IR) OpenDevice() bool {
	Set_Debug(i.CTX, 3)

	detect := false
	devices, _ := Get_Device_List(i.CTX)
	for _, device := range devices {
		fd, _ := Get_Device_Descriptor(device)
		if fd.IdVendor == VENDOR_ID && fd.IdProduct == PRODUCT_ID {
			detect = true
		}
	}

	if !detect {
		fmt.Println("device not found")
		return false
	}

	dh := Open_Device_With_VID_PID(i.CTX, VENDOR_ID, PRODUCT_ID)

	if ok, err := Kernel_Driver_Active(dh, 3); ok {
		if err = Detach_Kernel_Driver(dh, 3); err != nil {
			fmt.Println(err)
			return false
		}
	}

	if err := Claim_Interface(dh, 3); err != nil {
		fmt.Println(err)
		return false
	}

	i.DH = dh
	return true
}

func (i *IR) Write(freq uint64, data []byte) error {
	data_count := uint32(len(data))

	// 送信ビット数　リーダーコード + コード + 終了コード
	send_bit_num := (data_count / 4) + ((data_count % 4) / 2)

	if (i.DH == nil) &&
		(freq < IR_FREQ_MIN && IR_FREQ_MAX < freq) &&
		(send_bit_num < 1) &&
		(IR_SEND_DATA_MAX_LEN_K < send_bit_num) {
		return fmt.Errorf("invalid argument")
	}

	outbuffer := make([]byte, 64)
	inbuffer := make([]byte, 64)
	var send_bit_pos uint32 = 0
	var set_bit_size uint32 = 0
	for index := 0; ; {
		index++
		outbuffer[0] = 0x34
		outbuffer[1] = (byte)((send_bit_num >> 8) & 0xFF)
		outbuffer[2] = (byte)(send_bit_num & 0xFF)
		outbuffer[3] = (byte)((send_bit_pos >> 8) & 0xFF)
		outbuffer[4] = (byte)(send_bit_pos & 0xFF)

		if send_bit_num > send_bit_pos {
			set_bit_size = send_bit_num - send_bit_pos
			if set_bit_size > IR_SEND_DATA_USB_SEND_MAX_LEN {
				set_bit_size = IR_SEND_DATA_USB_SEND_MAX_LEN
			}
		} else {
			set_bit_size = 0
		}

		outbuffer[5] = (byte)(set_bit_size & 0xFF)

		if set_bit_size <= 0 {
			break
		}

		// set ir data
		var fi uint32 = 0
		for fi = 0; fi < set_bit_size; fi++ {
			// ON Count
			outbuffer[6+(fi*4)] = data[send_bit_pos*4]
			outbuffer[6+(fi*4)+1] = data[(send_bit_pos*4)+1]
			// OFF Count
			outbuffer[6+(fi*4)+2] = data[(send_bit_pos*4)+2]
			outbuffer[6+(fi*4)+3] = data[(send_bit_pos*4)+3]
			send_bit_pos++
		}

		_, err := Interrupt_Transfer(i.DH, BTO_EP_OUT, outbuffer, 5000)
		if err != nil {
			return err
		}
		//Now get the response packet from the firmware.
		_, err = Interrupt_Transfer(i.DH, BTO_EP_IN, inbuffer, 5000)
		if err != nil {
			return err
		}
		//INBuffer[0] is an echo back of the command (see microcontroller firmware).
		//INBuffer[1] contains the I/O port pin value for the pushbutton (see microcontroller firmware).
		if inbuffer[0] == 0x34 && inbuffer[1] != 0x00 {
			return fmt.Errorf("detaset error")
		}
	}

	// データ送信要求セット
	outbuffer[0] = 0x35 //0x81 is the "Get Pushbutton State" command in the firmware
	outbuffer[1] = (byte)((freq >> 8) & 0xFF)
	outbuffer[2] = (byte)(freq & 0xFF)
	outbuffer[3] = (byte)((send_bit_num >> 8) & 0xFF)
	outbuffer[4] = (byte)(send_bit_num & 0xFF)

	//To get the pushbutton state, first, we send a packet with our "Get Pushbutton State" command in it.
	_, err := Interrupt_Transfer(i.DH, BTO_EP_OUT, outbuffer, 5000)
	if err != nil {
		return err
	}

	//Now get the response packet from the firmware.
	_, err = Interrupt_Transfer(i.DH, BTO_EP_IN, inbuffer, 5000)
	if err != nil {
		return err
	}

	if inbuffer[0] == 0x35 && inbuffer[1] == 0x00 {
		return nil
	}

	return nil
}
