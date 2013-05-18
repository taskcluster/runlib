package win32

import (
	"syscall"
	"unsafe"
)

var (
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	psapi    = syscall.NewLazyDLL("psapi.dll")
	userenv  = syscall.NewLazyDLL("userenv.dll")
	user32   = syscall.NewLazyDLL("user32.dll")

	procCreateProcessWithLogonW   = advapi32.NewProc("CreateProcessWithLogonW")
	procCreateProcessAsUserW      = advapi32.NewProc("CreateProcessAsUserW")
	procResumeThread              = kernel32.NewProc("ResumeThread")
	procGetProcessMemoryInfo      = psapi.NewProc("GetProcessMemoryInfo")
	procLogonUserW                = advapi32.NewProc("LogonUserW")
	procLoadUserProfileW          = userenv.NewProc("LoadUserProfileW")
	procUnloadUserProfile         = userenv.NewProc("UnloadUserProfile")
	procGetProcessWindowStation   = user32.NewProc("GetProcessWindowStation")
	procGetCurrentThreadId        = kernel32.NewProc("GetCurrentThreadId")
	procGetThreadDesktop          = user32.NewProc("GetThreadDesktop")
	procCreateWindowStationW      = user32.NewProc("CreateWindowStationW")
	procSetProcessWindowStation   = user32.NewProc("SetProcessWindowStation")
	procCreateDesktopW            = user32.NewProc("CreateDesktopW")
	procSetThreadDesktop          = user32.NewProc("SetThreadDesktop")
	procGetUserObjectInformationW = user32.NewProc("GetUserObjectInformationW")
	procCloseWindowStation        = user32.NewProc("CloseWindowStation")
	procCreateJobObjectW          = kernel32.NewProc("CreateJobObjectW")
	procQueryInformationJobObject = kernel32.NewProc("QueryInformationJobObject")
	procSetInformationJobObject   = kernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject  = kernel32.NewProc("AssignProcessToJobObject")
	procVirtualAllocEx            = kernel32.NewProc("VirtualAllocEx")
	procWriteProcessMemory        = kernel32.NewProc("WriteProcessMemory")
	procGetModuleHandleW          = kernel32.NewProc("GetModuleHandleW")
	procCreateRemoteThread        = kernel32.NewProc("CreateRemoteThread")
	procVirtualFreeEx             = kernel32.NewProc("VirtualFreeEx")
)

const (
	CREATE_BREAKAWAY_FROM_JOB = 0x01000000
	CREATE_NEW_CONSOLE        = 0x00000010
	CREATE_NEW_PROCESS_GROUP  = 0x00000200
	CREATE_SUSPENDED          = 0x00000004

	LOGON_WITH_PROFILE = 0x00000001

	STARTF_FORCEOFFFEEDBACK = 0x00000080

	FILE_FLAG_SEQUENTIAL_SCAN = 0x08000000

	LOGON32_PROVIDER_DEFAULT = 0
	LOGON32_PROVIDER_WINNT35 = 1
	LOGON32_PROVIDER_WINNT40 = 2
	LOGON32_PROVIDER_WINNT50 = 3

	LOGON32_LOGON_INTERACTIVE       = 2
	LOGON32_LOGON_NETWORK           = 3
	LOGON32_LOGON_BATCH             = 4
	LOGON32_LOGON_SERVICE           = 5
	LOGON32_LOGON_UNLOCK            = 7
	LOGON32_LOGON_NETWORK_CLEARTEXT = 8
	LOGON32_LOGON_NEW_CREDENTIALS   = 9

	MAXIMUM_ALLOWED = 0x2000000
	PI_NOUI         = 2
)

type ProcessMemoryCountersEx struct {
	Cb                         uint32  // DWORD
	PageFaultCount             uint32  // DWORD
	PeakWorkingSetSize         uintptr // SIZE_T
	WorkingSetSize             uintptr // SIZE_T
	QuotaPeakPagedPoolUsage    uintptr // SIZE_T
	QuotaPagedPoolUsage        uintptr // SIZE_T
	QuotaPeakNonPagedPoolUsage uintptr // SIZE_T
	QuotaNonPagedPoolUsage     uintptr // SIZE_T
	PagefileUsage              uintptr // SIZE_T
	PeakPagefileUsage          uintptr // SIZE_T
	PrivateUsage               uintptr // SIZE_T
}

type ProfileInfo struct {
	Size        uint32
	Flags       uint32
	Username    *uint16
	ProfilePath *uint16
	DefaultPath *uint16
	ServerName  *uint16
	PolicyPath  *uint16
	Profile     syscall.Handle
}

type Hwinsta uintptr
type Hdesk uintptr

func MakeInheritSa() *syscall.SecurityAttributes {
	var sa syscall.SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
}

func StringPtrToUTF16Ptr(src *string) (result *uint16) {
	if src != nil {
		return syscall.StringToUTF16Ptr(*src)
	}
	return nil
}

func ListToEnvironmentBlock(list *[]string) *uint16 {
	if list == nil {
		return nil
	}

	size := 1
	for _, v := range *list {
		size += len(syscall.StringToUTF16(v))
	}

	result := make([]uint16, size)

	tail := 0

	for _, v := range *list {
		uline := syscall.StringToUTF16(v)
		copy(result[tail:], uline)
		tail += len(uline)
	}

	result[tail] = 0

	return &result[0]
}

func CreateProcessWithLogonW(
	username *uint16,
	domain *uint16,
	password *uint16,
	logonFlags uint32,
	applicationName *uint16,
	commandLine *uint16,
	creationFlags uint32,
	environment *uint16,
	currentDirectory *uint16,
	startupInfo *syscall.StartupInfo,
	processInformation *syscall.ProcessInformation) (err error) {

	r1, _, e1 := procCreateProcessWithLogonW.Call(
		uintptr(unsafe.Pointer(username)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(password)),
		uintptr(logonFlags),
		uintptr(unsafe.Pointer(applicationName)),
		uintptr(unsafe.Pointer(commandLine)),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(environment)), // env
		uintptr(unsafe.Pointer(currentDirectory)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(processInformation)))

	if int(r1) == 0 {
		return e1
	}

	return nil
}

func CreateProcessAsUser(
	token syscall.Handle,
	applicationName *uint16,
	commandLine *uint16,
	procSecurity *syscall.SecurityAttributes,
	threadSecurity *syscall.SecurityAttributes,
	inheritHandles bool,
	creationFlags uint32,
	environment *uint16,
	currentDirectory *uint16,
	startupInfo *syscall.StartupInfo,
	processInformation *syscall.ProcessInformation) (err error) {

	var _p0 uint32
	if inheritHandles {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r1, _, e1 := procCreateProcessAsUserW.Call(
		uintptr(token),
		uintptr(unsafe.Pointer(applicationName)),
		uintptr(unsafe.Pointer(commandLine)),
		uintptr(unsafe.Pointer(procSecurity)),
		uintptr(unsafe.Pointer(threadSecurity)),
		uintptr(_p0),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(environment)), // env
		uintptr(unsafe.Pointer(currentDirectory)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(processInformation)))

	if int(r1) == 0 {
		return e1
	}

	return nil
}

func ResumeThread(thread syscall.Handle) (suspendCount int, err error) {
	r1, _, e1 := procResumeThread.Call(uintptr(thread))
	if int(r1) == -1 {
		return -1, e1
	}
	return int(r1), nil
}

func GetProcessMemoryInfo(process syscall.Handle) (pmc *ProcessMemoryCountersEx, err error) {
	pmc = &ProcessMemoryCountersEx{}
	pmc.Cb = uint32(unsafe.Sizeof(*pmc))
	r1, _, e1 := procGetProcessMemoryInfo.Call(uintptr(process), uintptr(unsafe.Pointer(pmc)), uintptr(pmc.Cb))
	if int(r1) == 0 {
		return nil, e1
	}
	return pmc, nil
}

func LogonUser(username *uint16, domain *uint16, password *uint16, logonType uint32, logonProvider uint32) (token syscall.Handle, err error) {
	r1, _, e1 := procLogonUserW.Call(
		uintptr(unsafe.Pointer(username)),
		uintptr(unsafe.Pointer(domain)),
		uintptr(unsafe.Pointer(password)),
		uintptr(logonType),
		uintptr(logonProvider),
		uintptr(unsafe.Pointer(&token)))

	if int(r1) == 0 {
		return syscall.InvalidHandle, e1
	}
	return
}

func LoadUserProfile(token syscall.Handle, pinfo *ProfileInfo) error {
	r1, _, e1 := procLoadUserProfileW.Call(
		uintptr(token),
		uintptr(unsafe.Pointer(pinfo)))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

func UnloadUserProfile(token, profile syscall.Handle) error {
	r1, _, e1 := procUnloadUserProfile.Call(
		uintptr(token),
		uintptr(profile))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

func GetProcessWindowStation() (Hwinsta, error) {
	r1, _, e1 := procGetProcessWindowStation.Call()
	if int(r1) == 0 {
		return Hwinsta(r1), e1
	}
	return Hwinsta(r1), nil
}

func GetCurrentThreadId() uint32 {
	r1, _, _ := procGetCurrentThreadId.Call()
	return uint32(r1)
}

func GetThreadDesktop(threadId uint32) (Hdesk, error) {
	r1, _, e1 := procGetThreadDesktop.Call(
		uintptr(threadId))
	if int(r1) == 0 {
		return Hdesk(r1), e1
	}
	return Hdesk(r1), nil
}

func CreateWindowStation(winsta *uint16, flags, desiredAccess uint32, sa *syscall.SecurityAttributes) (Hwinsta, error) {
	r1, _, e1 := procCreateWindowStationW.Call(
		uintptr(unsafe.Pointer(winsta)),
		uintptr(flags),
		uintptr(desiredAccess),
		uintptr(unsafe.Pointer(sa)))
	if int(r1) == 0 {
		return Hwinsta(r1), e1
	}
	return Hwinsta(r1), nil
}

func SetProcessWindowStation(winsta Hwinsta) error {
	r1, _, e1 := procSetProcessWindowStation.Call(
		uintptr(winsta))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

func CreateDesktop(desktop, device *uint16, devmode uintptr, flags, desiredAccess uint32, sa *syscall.SecurityAttributes) (Hdesk, error) {
	r1, _, e1 := procCreateDesktopW.Call(
		uintptr(unsafe.Pointer(desktop)),
		uintptr(unsafe.Pointer(device)),
		devmode,
		uintptr(flags),
		uintptr(desiredAccess),
		uintptr(unsafe.Pointer(sa)))
	if int(r1) == 0 {
		return Hdesk(r1), e1
	}
	return Hdesk(r1), nil
}

func SetThreadDesktop(desk Hdesk) error {
	r1, _, e1 := procSetThreadDesktop.Call(
		uintptr(desk))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

const (
	UOI_NAME = 2
)

func GetUserObjectInformation(obj syscall.Handle, index int, info unsafe.Pointer, length uint32) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procGetUserObjectInformationW.Call(
		uintptr(obj),
		uintptr(index),
		uintptr(info),
		uintptr(length),
		uintptr(unsafe.Pointer(&nLength)))
	if int(r1) == 0 {
		return nLength, e1
	}
	return 0, nil
}

func GetUserObjectName(obj syscall.Handle) (string, error) {
	namebuf := make([]uint16, 256)
	_, err := GetUserObjectInformation(obj, UOI_NAME, unsafe.Pointer(&namebuf[0]), 256*2)
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(namebuf), nil
}

func CloseWindowStation(winsta Hwinsta) error {
	r1, _, e1 := procCloseWindowStation.Call(
		uintptr(winsta))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

func CreateJobObject(sa *syscall.SecurityAttributes, name *uint16) (syscall.Handle, error) {
	r1, _, e1 := procCreateJobObjectW.Call(
		uintptr(unsafe.Pointer(sa)),
		uintptr(unsafe.Pointer(name)))
	if int(r1) == 0 {
		return syscall.InvalidHandle, e1
	}
	return syscall.Handle(r1), nil
}

func QueryInformationJobObject(job syscall.Handle, infoclass uint32, info unsafe.Pointer, length uint32) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procQueryInformationJobObject.Call(
		uintptr(job),
		uintptr(infoclass),
		uintptr(info),
		uintptr(length),
		uintptr(unsafe.Pointer(&nLength)))

	if int(r1) == 0 {
		return nLength, e1
	}
	return nLength, nil
}

func SetInformationJobObject(job syscall.Handle, infoclass uint32, info unsafe.Pointer, length uint32) error {
	r1, _, e1 := procSetInformationJobObject.Call(
		uintptr(job),
		uintptr(infoclass),
		uintptr(info),
		uintptr(length))

	if int(r1) == 0 {
		return e1
	}
	return nil
}

type JobObjectBasicAccountingInformation struct {
	TotalUserTime             uint64
	TotalKernelTime           uint64
	ThisPeriodTotalUserTime   uint64
	ThisPeriodTotalKernelTime uint64
	TotalPageFaultCount       uint32
	TotalProcesses            uint32
	ActiveProcesses           uint32
	TotalTerminatedProcesses  uint32
}

type JobObjectBasicUiRestrictions struct {
	UIRestrictionClass uint32
}

const (
	JOB_OBJECT_UILIMIT_DESKTOP          = 0x40
	JOB_OBJECT_UILIMIT_DISPLAYSETTINGS  = 0x10
	JOB_OBJECT_UILIMIT_EXITWINDOWS      = 0x80
	JOB_OBJECT_UILIMIT_GLOBALATOMS      = 0x20
	JOB_OBJECT_UILIMIT_HANDLES          = 1
	JOB_OBJECT_UILIMIT_READCLIPBOARD    = 2
	JOB_OBJECT_UILIMIT_SYSTEMPARAMETERS = 8
	JOB_OBJECT_UILIMIT_WRITECLIPBOARD   = 4
)

func GetJobObjectBasicAccountingInformation(job syscall.Handle) (*JobObjectBasicAccountingInformation, error) {
	var jinfo JobObjectBasicAccountingInformation
	_, err := QueryInformationJobObject(job, 1, unsafe.Pointer(&jinfo), uint32(unsafe.Sizeof(jinfo)))
	if err != nil {
		return nil, err
	}
	return &jinfo, nil
}

type JobObjectBasicLimitInformation struct {
	PerProcessUserTimeLimit uint64  // LARGE_INTEGER
	PerJobUserTimeLimit     uint64  // LARGE_INTEGER
	LimitFlags              uint32  // DWORD
	MinimumWorkingSetSize   uintptr // SIZE_T
	MaximumWorkingSetSize   uintptr // SIZE_T
	ActiveProcessLimit      uint32  // DWORD
	Affinity                uintptr // ULONG_PTR
	PriorityClass           uint32  // DWORD
	SchedulingClass         uint32  // DWORD
}

const (
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x2000
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x400
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 8
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x200
	JOB_OBJECT_LIMIT_JOB_TIME                   = 4
	JOB_OBJECT_LIMIT_PROCESS_MEMORY             = 0x100
	JOB_OBJECT_LIMIT_PROCESS_TIME               = 2
	JOB_OBJECT_LIMIT_WORKINGSET                 = 1
)

type IoCounters struct {
	ReadOperationCount  uint64 // ULONGLONG
	WriteOperationCount uint64 // ULONGLONG
	OtherOperationCount uint64 // ULONGLONG
	ReadTransferCount   uint64 // ULONGLONG
	WriteTransferCount  uint64 // ULONGLONG
	OtherTransferCount  uint64 // ULONGLONG
}


func GetJobObjectExtendedLimitInformation(job syscall.Handle) (*JobObjectExtendedLimitInformation, error) {
	var jinfo JobObjectExtendedLimitInformation
	_, err := QueryInformationJobObject(job, 9, unsafe.Pointer(&jinfo), uint32(unsafe.Sizeof(jinfo)))
	if err != nil {
		return nil, err
	}
	return &jinfo, nil
}

func SetJobObjectBasicUiRestrictions(job syscall.Handle, info *JobObjectBasicUiRestrictions) error {
	return SetInformationJobObject(job, 4, unsafe.Pointer(info), uint32(unsafe.Sizeof(*info)))
}

func SetJobObjectExtendedLimitInformation(job syscall.Handle, info *JobObjectExtendedLimitInformation) error {
	return SetInformationJobObject(job, 9, unsafe.Pointer(info), uint32(unsafe.Sizeof(*info)))
}

func AssignProcessToJobObject(job syscall.Handle, process syscall.Handle) error {
	r1, _, e1 := procAssignProcessToJobObject.Call(
		uintptr(job),
		uintptr(process))
	if int(r1) == 0 {
		return e1
	}
	return nil
}

const (
	MEM_COMMIT     = 0x00001000
	PAGE_READWRITE = 0x04
)

func VirtualAllocEx(process syscall.Handle, addr uintptr, size, allocType, protect uint32) (uintptr, error) {
	r1, _, e1 := procVirtualAllocEx.Call(
		uintptr(process),
		addr,
		uintptr(size),
		uintptr(allocType),
		uintptr(protect))

	if int(r1) == 0 {
		return r1, e1
	}
	return r1, nil
}

func WriteProcessMemory(process syscall.Handle, addr uintptr, buf unsafe.Pointer, size uint32) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procWriteProcessMemory.Call(
		uintptr(process),
		addr,
		uintptr(buf),
		uintptr(size),
		uintptr(unsafe.Pointer(&nLength)))

	if int(r1) == 0 {
		return nLength, e1
	}
	return nLength, nil
}

func GetModuleHandle(name *uint16) (syscall.Handle, error) {
	r1, _, e1 := procGetModuleHandleW.Call(uintptr(unsafe.Pointer(name)))
	if int(r1) == 0 {
		return syscall.InvalidHandle, e1
	}
	return syscall.Handle(r1), nil
}

func CreateRemoteThread(process syscall.Handle, sa *syscall.SecurityAttributes, stackSize uint32, startAddress,
	parameter uintptr, creationFlags uint32) (syscall.Handle, uint32, error) {
	var threadId uint32
	r1, _, e1 := procCreateRemoteThread.Call(
		uintptr(process),
		uintptr(unsafe.Pointer(sa)),
		uintptr(stackSize),
		startAddress,
		parameter,
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(&threadId)))

	if int(r1) == 0 {
		return syscall.InvalidHandle, 0, e1
	}
	return syscall.Handle(r1), threadId, nil
}

const (
	MEM_RELEASE = 0x8000
)

func VirtualFreeEx(process syscall.Handle, addr uintptr, size, freeType uint32) error {
	r1, _, e1 := procVirtualFreeEx.Call(
		uintptr(process),
		addr,
		uintptr(size),
		uintptr(freeType))

	if int(r1) == 0 {
		return e1
	}
	return nil
}

func SetInheritHandle(h syscall.Handle, inherit bool) error {
	var v uint32
	if inherit {
		v = syscall.HANDLE_FLAG_INHERIT
	}
	return syscall.SetHandleInformation(syscall.Handle(h), syscall.HANDLE_FLAG_INHERIT, v)
}