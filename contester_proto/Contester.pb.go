// Code generated by protoc-gen-gogo.
// source: Contester.proto
// DO NOT EDIT!

package contester_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Compilation_Code int32

const (
	Compilation_Success Compilation_Code = 1
	Compilation_Failure Compilation_Code = 2
)

var Compilation_Code_name = map[int32]string{
	1: "Success",
	2: "Failure",
}
var Compilation_Code_value = map[string]int32{
	"Success": 1,
	"Failure": 2,
}

func (x Compilation_Code) Enum() *Compilation_Code {
	p := new(Compilation_Code)
	*p = x
	return p
}
func (x Compilation_Code) String() string {
	return proto.EnumName(Compilation_Code_name, int32(x))
}
func (x *Compilation_Code) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Compilation_Code_value, data, "Compilation_Code")
	if err != nil {
		return err
	}
	*x = Compilation_Code(value)
	return nil
}
func (Compilation_Code) EnumDescriptor() ([]byte, []int) { return fileDescriptorContester, []int{0, 0} }

type Compilation struct {
	Failure          *bool                 `protobuf:"varint,1,opt,name=failure" json:"failure,omitempty"`
	ResultSteps      []*Compilation_Result `protobuf:"bytes,2,rep,name=result_steps,json=resultSteps" json:"result_steps,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *Compilation) Reset()                    { *m = Compilation{} }
func (m *Compilation) String() string            { return proto.CompactTextString(m) }
func (*Compilation) ProtoMessage()               {}
func (*Compilation) Descriptor() ([]byte, []int) { return fileDescriptorContester, []int{0} }

func (m *Compilation) GetFailure() bool {
	if m != nil && m.Failure != nil {
		return *m.Failure
	}
	return false
}

func (m *Compilation) GetResultSteps() []*Compilation_Result {
	if m != nil {
		return m.ResultSteps
	}
	return nil
}

type Compilation_Result struct {
	StepName         *string         `protobuf:"bytes,1,opt,name=step_name,json=stepName" json:"step_name,omitempty"`
	Execution        *LocalExecution `protobuf:"bytes,2,opt,name=execution" json:"execution,omitempty"`
	Failure          *bool           `protobuf:"varint,3,opt,name=failure" json:"failure,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *Compilation_Result) Reset()                    { *m = Compilation_Result{} }
func (m *Compilation_Result) String() string            { return proto.CompactTextString(m) }
func (*Compilation_Result) ProtoMessage()               {}
func (*Compilation_Result) Descriptor() ([]byte, []int) { return fileDescriptorContester, []int{0, 0} }

func (m *Compilation_Result) GetStepName() string {
	if m != nil && m.StepName != nil {
		return *m.StepName
	}
	return ""
}

func (m *Compilation_Result) GetExecution() *LocalExecution {
	if m != nil {
		return m.Execution
	}
	return nil
}

func (m *Compilation_Result) GetFailure() bool {
	if m != nil && m.Failure != nil {
		return *m.Failure
	}
	return false
}

func init() {
	proto.RegisterType((*Compilation)(nil), "contester.proto.Compilation")
	proto.RegisterType((*Compilation_Result)(nil), "contester.proto.Compilation.Result")
	proto.RegisterEnum("contester.proto.Compilation_Code", Compilation_Code_name, Compilation_Code_value)
}
func (m *Compilation) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Compilation) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Failure != nil {
		data[i] = 0x8
		i++
		if *m.Failure {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if len(m.ResultSteps) > 0 {
		for _, msg := range m.ResultSteps {
			data[i] = 0x12
			i++
			i = encodeVarintContester(data, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(data[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Compilation_Result) Marshal() (data []byte, err error) {
	size := m.Size()
	data = make([]byte, size)
	n, err := m.MarshalTo(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

func (m *Compilation_Result) MarshalTo(data []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.StepName != nil {
		data[i] = 0xa
		i++
		i = encodeVarintContester(data, i, uint64(len(*m.StepName)))
		i += copy(data[i:], *m.StepName)
	}
	if m.Execution != nil {
		data[i] = 0x12
		i++
		i = encodeVarintContester(data, i, uint64(m.Execution.Size()))
		n1, err := m.Execution.MarshalTo(data[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Failure != nil {
		data[i] = 0x18
		i++
		if *m.Failure {
			data[i] = 1
		} else {
			data[i] = 0
		}
		i++
	}
	if m.XXX_unrecognized != nil {
		i += copy(data[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeFixed64Contester(data []byte, offset int, v uint64) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	data[offset+4] = uint8(v >> 32)
	data[offset+5] = uint8(v >> 40)
	data[offset+6] = uint8(v >> 48)
	data[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Contester(data []byte, offset int, v uint32) int {
	data[offset] = uint8(v)
	data[offset+1] = uint8(v >> 8)
	data[offset+2] = uint8(v >> 16)
	data[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintContester(data []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		data[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	data[offset] = uint8(v)
	return offset + 1
}
func (m *Compilation) Size() (n int) {
	var l int
	_ = l
	if m.Failure != nil {
		n += 2
	}
	if len(m.ResultSteps) > 0 {
		for _, e := range m.ResultSteps {
			l = e.Size()
			n += 1 + l + sovContester(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Compilation_Result) Size() (n int) {
	var l int
	_ = l
	if m.StepName != nil {
		l = len(*m.StepName)
		n += 1 + l + sovContester(uint64(l))
	}
	if m.Execution != nil {
		l = m.Execution.Size()
		n += 1 + l + sovContester(uint64(l))
	}
	if m.Failure != nil {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovContester(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozContester(x uint64) (n int) {
	return sovContester(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Compilation) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContester
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Compilation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Compilation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Failure", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContester
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Failure = &b
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResultSteps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContester
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthContester
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ResultSteps = append(m.ResultSteps, &Compilation_Result{})
			if err := m.ResultSteps[len(m.ResultSteps)-1].Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContester(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContester
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Compilation_Result) Unmarshal(data []byte) error {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContester
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Result: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Result: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StepName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContester
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContester
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(data[iNdEx:postIndex])
			m.StepName = &s
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Execution", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContester
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthContester
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Execution == nil {
				m.Execution = &LocalExecution{}
			}
			if err := m.Execution.Unmarshal(data[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Failure", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContester
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Failure = &b
		default:
			iNdEx = preIndex
			skippy, err := skipContester(data[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContester
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, data[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipContester(data []byte) (n int, err error) {
	l := len(data)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContester
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := data[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContester
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if data[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContester
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := data[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthContester
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowContester
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := data[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipContester(data[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthContester = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContester   = fmt.Errorf("proto: integer overflow")
)

var fileDescriptorContester = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x77, 0xce, 0xcf, 0x2b,
	0x49, 0x2d, 0x2e, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4f, 0x46, 0x15,
	0x90, 0xe2, 0xf6, 0xc9, 0x4f, 0x4e, 0xcc, 0x81, 0x70, 0x94, 0x26, 0x32, 0x71, 0x71, 0x3b, 0xe7,
	0xe7, 0x16, 0x64, 0xe6, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x09, 0x49, 0x70, 0xb1, 0xa7, 0x25, 0x66,
	0xe6, 0x94, 0x16, 0xa5, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x04, 0xc1, 0xb8, 0x42, 0x6e, 0x5c,
	0x3c, 0x45, 0xa9, 0xc5, 0xa5, 0x39, 0x25, 0xf1, 0x40, 0xb3, 0x0a, 0x8a, 0x25, 0x98, 0x14, 0x98,
	0x35, 0xb8, 0x8d, 0x94, 0xf5, 0xd0, 0x8c, 0xd7, 0x43, 0x32, 0x4d, 0x2f, 0x08, 0xac, 0x21, 0x88,
	0x1b, 0xa2, 0x31, 0x18, 0xa4, 0x4f, 0xaa, 0x8e, 0x8b, 0x0d, 0x22, 0x2c, 0x24, 0xcd, 0xc5, 0x09,
	0x32, 0x2a, 0x3e, 0x2f, 0x31, 0x17, 0x62, 0x1b, 0x67, 0x10, 0x07, 0x48, 0xc0, 0x0f, 0xc8, 0x17,
	0xb2, 0xe5, 0xe2, 0x4c, 0xad, 0x48, 0x4d, 0x2e, 0x05, 0x99, 0x03, 0xb4, 0x8b, 0x11, 0x68, 0x97,
	0x3c, 0x86, 0x5d, 0x60, 0x9f, 0xb8, 0xc2, 0x94, 0x05, 0x21, 0x74, 0x20, 0xfb, 0x83, 0x19, 0xc5,
	0x1f, 0x4a, 0x0a, 0x5c, 0x2c, 0xce, 0xf9, 0x29, 0xa9, 0x42, 0xdc, 0x5c, 0xec, 0xc1, 0xa5, 0xc9,
	0xc9, 0xa9, 0xc5, 0xc5, 0x02, 0x8c, 0x20, 0x8e, 0x1b, 0x44, 0x5e, 0x80, 0xc9, 0x49, 0xef, 0xc4,
	0x23, 0x39, 0xc6, 0x0b, 0x40, 0xfc, 0x00, 0x88, 0x67, 0x3c, 0x96, 0x63, 0xe0, 0x92, 0xc9, 0x2f,
	0x4a, 0xd7, 0x2b, 0x2e, 0xc9, 0xcc, 0x4b, 0x2f, 0x4a, 0xac, 0x44, 0x77, 0x05, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x82, 0x0a, 0x10, 0xf2, 0x73, 0x01, 0x00, 0x00,
}
