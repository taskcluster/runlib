package win32

// Refer to
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa383751(v=vs.85).aspx
// for understanding the c++ -> go type mappings

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/taskcluster/ntr"

	"golang.org/x/text/encoding/charmap"
)

var (
	shell32 = NewLazyDLL("shell32.dll")
	ole32   = NewLazyDLL("ole32.dll")
	secur32 = NewLazyDLL("secur32.dll")

	procCloseDesktop                   = user32.NewProc("CloseDesktop")
	procSwitchDesktop                  = user32.NewProc("SwitchDesktop")
	procSetPriorityClass               = kernel32.NewProc("SetPriorityClass")
	procCreateEnvironmentBlock         = userenv.NewProc("CreateEnvironmentBlock")
	procDestroyEnvironmentBlock        = userenv.NewProc("DestroyEnvironmentBlock")
	procSHSetKnownFolderPath           = shell32.NewProc("SHSetKnownFolderPath")
	procSHGetKnownFolderPath           = shell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree                  = ole32.NewProc("CoTaskMemFree")
	procCloseHandle                    = kernel32.NewProc("CloseHandle")
	procLsaConnectUntrusted            = secur32.NewProc("LsaConnectUntrusted")
	procLsaLookupAuthenticationPackage = secur32.NewProc("LsaLookupAuthenticationPackage")
	procAllocateLocallyUniqueId        = advapi32.NewProc("AllocateLocallyUniqueId")
	procLsaLogonUser                   = secur32.NewProc("LsaLogonUser")
	procLsaCallAuthenticationPackage   = secur32.NewProc("LsaCallAuthenticationPackage")
	procLsaFreeReturnBuffer            = secur32.NewProc("LsaFreeReturnBuffer")
	procLsaRegisterLogonProcess        = secur32.NewProc("LsaRegisterLogonProcess")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/dd378457(v=vs.85).aspx
var (
	FOLDERID_AccountPictures        = syscall.GUID{Data1: 0x008CA0B1, Data2: 0x55B4, Data3: 0x4C56, Data4: [8]byte{0xB8, 0xA8, 0x4D, 0xE4, 0xB2, 0x99, 0xD3, 0xBE}}
	FOLDERID_AddNewPrograms         = syscall.GUID{Data1: 0xDE61D971, Data2: 0x5EBC, Data3: 0x4F02, Data4: [8]byte{0xA3, 0xA9, 0x6C, 0x82, 0x89, 0x5E, 0x5C, 0x04}}
	FOLDERID_AdminTools             = syscall.GUID{Data1: 0x724EF170, Data2: 0xA42D, Data3: 0x4FEF, Data4: [8]byte{0x9F, 0x26, 0xB6, 0x0E, 0x84, 0x6F, 0xBA, 0x4F}}
	FOLDERID_ApplicationShortcuts   = syscall.GUID{Data1: 0xA3918781, Data2: 0xE5F2, Data3: 0x4890, Data4: [8]byte{0xB3, 0xD9, 0xA7, 0xE5, 0x43, 0x32, 0x32, 0x8C}}
	FOLDERID_AppsFolder             = syscall.GUID{Data1: 0x1E87508D, Data2: 0x89C2, Data3: 0x42F0, Data4: [8]byte{0x8A, 0x7E, 0x64, 0x5A, 0x0F, 0x50, 0xCA, 0x58}}
	FOLDERID_AppUpdates             = syscall.GUID{Data1: 0xA305CE99, Data2: 0xF527, Data3: 0x492B, Data4: [8]byte{0x8B, 0x1A, 0x7E, 0x76, 0xFA, 0x98, 0xD6, 0xE4}}
	FOLDERID_CameraRoll             = syscall.GUID{Data1: 0xAB5FB87B, Data2: 0x7CE2, Data3: 0x4F83, Data4: [8]byte{0x91, 0x5D, 0x55, 0x08, 0x46, 0xC9, 0x53, 0x7B}}
	FOLDERID_CDBurning              = syscall.GUID{Data1: 0x9E52AB10, Data2: 0xF80D, Data3: 0x49DF, Data4: [8]byte{0xAC, 0xB8, 0x43, 0x30, 0xF5, 0x68, 0x78, 0x55}}
	FOLDERID_ChangeRemovePrograms   = syscall.GUID{Data1: 0xDF7266AC, Data2: 0x9274, Data3: 0x4867, Data4: [8]byte{0x8D, 0x55, 0x3B, 0xD6, 0x61, 0xDE, 0x87, 0x2D}}
	FOLDERID_CommonAdminTools       = syscall.GUID{Data1: 0xD0384E7D, Data2: 0xBAC3, Data3: 0x4797, Data4: [8]byte{0x8F, 0x14, 0xCB, 0xA2, 0x29, 0xB3, 0x92, 0xB5}}
	FOLDERID_CommonOEMLinks         = syscall.GUID{Data1: 0xC1BAE2D0, Data2: 0x10DF, Data3: 0x4334, Data4: [8]byte{0xBE, 0xDD, 0x7A, 0xA2, 0x0B, 0x22, 0x7A, 0x9D}}
	FOLDERID_CommonPrograms         = syscall.GUID{Data1: 0x0139D44E, Data2: 0x6AFE, Data3: 0x49F2, Data4: [8]byte{0x86, 0x90, 0x3D, 0xAF, 0xCA, 0xE6, 0xFF, 0xB8}}
	FOLDERID_CommonStartMenu        = syscall.GUID{Data1: 0xA4115719, Data2: 0xD62E, Data3: 0x491D, Data4: [8]byte{0xAA, 0x7C, 0xE7, 0x4B, 0x8B, 0xE3, 0xB0, 0x67}}
	FOLDERID_CommonStartup          = syscall.GUID{Data1: 0x82A5EA35, Data2: 0xD9CD, Data3: 0x47C5, Data4: [8]byte{0x96, 0x29, 0xE1, 0x5D, 0x2F, 0x71, 0x4E, 0x6E}}
	FOLDERID_CommonTemplates        = syscall.GUID{Data1: 0xB94237E7, Data2: 0x57AC, Data3: 0x4347, Data4: [8]byte{0x91, 0x51, 0xB0, 0x8C, 0x6C, 0x32, 0xD1, 0xF7}}
	FOLDERID_ComputerFolder         = syscall.GUID{Data1: 0x0AC0837C, Data2: 0xBBF8, Data3: 0x452A, Data4: [8]byte{0x85, 0x0D, 0x79, 0xD0, 0x8E, 0x66, 0x7C, 0xA7}}
	FOLDERID_ConflictFolder         = syscall.GUID{Data1: 0x4BFEFB45, Data2: 0x347D, Data3: 0x4006, Data4: [8]byte{0xA5, 0xBE, 0xAC, 0x0C, 0xB0, 0x56, 0x71, 0x92}}
	FOLDERID_ConnectionsFolder      = syscall.GUID{Data1: 0x6F0CD92B, Data2: 0x2E97, Data3: 0x45D1, Data4: [8]byte{0x88, 0xFF, 0xB0, 0xD1, 0x86, 0xB8, 0xDE, 0xDD}}
	FOLDERID_Contacts               = syscall.GUID{Data1: 0x56784854, Data2: 0xC6CB, Data3: 0x462B, Data4: [8]byte{0x81, 0x69, 0x88, 0xE3, 0x50, 0xAC, 0xB8, 0x82}}
	FOLDERID_ControlPanelFolder     = syscall.GUID{Data1: 0x82A74AEB, Data2: 0xAEB4, Data3: 0x465C, Data4: [8]byte{0xA0, 0x14, 0xD0, 0x97, 0xEE, 0x34, 0x6D, 0x63}}
	FOLDERID_Cookies                = syscall.GUID{Data1: 0x2B0F765D, Data2: 0xC0E9, Data3: 0x4171, Data4: [8]byte{0x90, 0x8E, 0x08, 0xA6, 0x11, 0xB8, 0x4F, 0xF6}}
	FOLDERID_Desktop                = syscall.GUID{Data1: 0xB4BFCC3A, Data2: 0xDB2C, Data3: 0x424C, Data4: [8]byte{0xB0, 0x29, 0x7F, 0xE9, 0x9A, 0x87, 0xC6, 0x41}}
	FOLDERID_DeviceMetadataStore    = syscall.GUID{Data1: 0x5CE4A5E9, Data2: 0xE4EB, Data3: 0x479D, Data4: [8]byte{0xB8, 0x9F, 0x13, 0x0C, 0x02, 0x88, 0x61, 0x55}}
	FOLDERID_Documents              = syscall.GUID{Data1: 0xFDD39AD0, Data2: 0x238F, Data3: 0x46AF, Data4: [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}
	FOLDERID_DocumentsLibrary       = syscall.GUID{Data1: 0x7B0DB17D, Data2: 0x9CD2, Data3: 0x4A93, Data4: [8]byte{0x97, 0x33, 0x46, 0xCC, 0x89, 0x02, 0x2E, 0x7C}}
	FOLDERID_Downloads              = syscall.GUID{Data1: 0x374DE290, Data2: 0x123F, Data3: 0x4565, Data4: [8]byte{0x91, 0x64, 0x39, 0xC4, 0x92, 0x5E, 0x46, 0x7B}}
	FOLDERID_Favorites              = syscall.GUID{Data1: 0x1777F761, Data2: 0x68AD, Data3: 0x4D8A, Data4: [8]byte{0x87, 0xBD, 0x30, 0xB7, 0x59, 0xFA, 0x33, 0xDD}}
	FOLDERID_Fonts                  = syscall.GUID{Data1: 0xFD228CB7, Data2: 0xAE11, Data3: 0x4AE3, Data4: [8]byte{0x86, 0x4C, 0x16, 0xF3, 0x91, 0x0A, 0xB8, 0xFE}}
	FOLDERID_Games                  = syscall.GUID{Data1: 0xCAC52C1A, Data2: 0xB53D, Data3: 0x4EDC, Data4: [8]byte{0x92, 0xD7, 0x6B, 0x2E, 0x8A, 0xC1, 0x94, 0x34}}
	FOLDERID_GameTasks              = syscall.GUID{Data1: 0x054FAE61, Data2: 0x4DD8, Data3: 0x4787, Data4: [8]byte{0x80, 0xB6, 0x09, 0x02, 0x20, 0xC4, 0xB7, 0x00}}
	FOLDERID_History                = syscall.GUID{Data1: 0xD9DC8A3B, Data2: 0xB784, Data3: 0x432E, Data4: [8]byte{0xA7, 0x81, 0x5A, 0x11, 0x30, 0xA7, 0x59, 0x63}}
	FOLDERID_HomeGroup              = syscall.GUID{Data1: 0x52528A6B, Data2: 0xB9E3, Data3: 0x4ADD, Data4: [8]byte{0xB6, 0x0D, 0x58, 0x8C, 0x2D, 0xBA, 0x84, 0x2D}}
	FOLDERID_HomeGroupCurrentUser   = syscall.GUID{Data1: 0x9B74B6A3, Data2: 0x0DFD, Data3: 0x4F11, Data4: [8]byte{0x9E, 0x78, 0x5F, 0x78, 0x00, 0xF2, 0xE7, 0x72}}
	FOLDERID_ImplicitAppShortcuts   = syscall.GUID{Data1: 0xBCB5256F, Data2: 0x79F6, Data3: 0x4CEE, Data4: [8]byte{0xB7, 0x25, 0xDC, 0x34, 0xE4, 0x02, 0xFD, 0x46}}
	FOLDERID_InternetCache          = syscall.GUID{Data1: 0x352481E8, Data2: 0x33BE, Data3: 0x4251, Data4: [8]byte{0xBA, 0x85, 0x60, 0x07, 0xCA, 0xED, 0xCF, 0x9D}}
	FOLDERID_InternetFolder         = syscall.GUID{Data1: 0x4D9F7874, Data2: 0x4E0C, Data3: 0x4904, Data4: [8]byte{0x96, 0x7B, 0x40, 0xB0, 0xD2, 0x0C, 0x3E, 0x4B}}
	FOLDERID_Libraries              = syscall.GUID{Data1: 0x1B3EA5DC, Data2: 0xB587, Data3: 0x4786, Data4: [8]byte{0xB4, 0xEF, 0xBD, 0x1D, 0xC3, 0x32, 0xAE, 0xAE}}
	FOLDERID_Links                  = syscall.GUID{Data1: 0xBFB9D5E0, Data2: 0xC6A9, Data3: 0x404C, Data4: [8]byte{0xB2, 0xB2, 0xAE, 0x6D, 0xB6, 0xAF, 0x49, 0x68}}
	FOLDERID_LocalAppData           = syscall.GUID{Data1: 0xF1B32785, Data2: 0x6FBA, Data3: 0x4FCF, Data4: [8]byte{0x9D, 0x55, 0x7B, 0x8E, 0x7F, 0x15, 0x70, 0x91}}
	FOLDERID_LocalAppDataLow        = syscall.GUID{Data1: 0xA520A1A4, Data2: 0x1780, Data3: 0x4FF6, Data4: [8]byte{0xBD, 0x18, 0x16, 0x73, 0x43, 0xC5, 0xAF, 0x16}}
	FOLDERID_LocalizedResourcesDir  = syscall.GUID{Data1: 0x2A00375E, Data2: 0x224C, Data3: 0x49DE, Data4: [8]byte{0xB8, 0xD1, 0x44, 0x0D, 0xF7, 0xEF, 0x3D, 0xDC}}
	FOLDERID_Music                  = syscall.GUID{Data1: 0x4BD8D571, Data2: 0x6D19, Data3: 0x48D3, Data4: [8]byte{0xBE, 0x97, 0x42, 0x22, 0x20, 0x08, 0x0E, 0x43}}
	FOLDERID_MusicLibrary           = syscall.GUID{Data1: 0x2112AB0A, Data2: 0xC86A, Data3: 0x4FFE, Data4: [8]byte{0xA3, 0x68, 0x0D, 0xE9, 0x6E, 0x47, 0x01, 0x2E}}
	FOLDERID_NetHood                = syscall.GUID{Data1: 0xC5ABBF53, Data2: 0xE17F, Data3: 0x4121, Data4: [8]byte{0x89, 0x00, 0x86, 0x62, 0x6F, 0xC2, 0xC9, 0x73}}
	FOLDERID_NetworkFolder          = syscall.GUID{Data1: 0xD20BEEC4, Data2: 0x5CA8, Data3: 0x4905, Data4: [8]byte{0xAE, 0x3B, 0xBF, 0x25, 0x1E, 0xA0, 0x9B, 0x53}}
	FOLDERID_OriginalImages         = syscall.GUID{Data1: 0x2C36C0AA, Data2: 0x5812, Data3: 0x4B87, Data4: [8]byte{0xBF, 0xD0, 0x4C, 0xD0, 0xDF, 0xB1, 0x9B, 0x39}}
	FOLDERID_PhotoAlbums            = syscall.GUID{Data1: 0x69D2CF90, Data2: 0xFC33, Data3: 0x4FB7, Data4: [8]byte{0x9A, 0x0C, 0xEB, 0xB0, 0xF0, 0xFC, 0xB4, 0x3C}}
	FOLDERID_PicturesLibrary        = syscall.GUID{Data1: 0xA990AE9F, Data2: 0xA03B, Data3: 0x4E80, Data4: [8]byte{0x94, 0xBC, 0x99, 0x12, 0xD7, 0x50, 0x41, 0x04}}
	FOLDERID_Pictures               = syscall.GUID{Data1: 0x33E28130, Data2: 0x4E1E, Data3: 0x4676, Data4: [8]byte{0x83, 0x5A, 0x98, 0x39, 0x5C, 0x3B, 0xC3, 0xBB}}
	FOLDERID_Playlists              = syscall.GUID{Data1: 0xDE92C1C7, Data2: 0x837F, Data3: 0x4F69, Data4: [8]byte{0xA3, 0xBB, 0x86, 0xE6, 0x31, 0x20, 0x4A, 0x23}}
	FOLDERID_PrintersFolder         = syscall.GUID{Data1: 0x76FC4E2D, Data2: 0xD6AD, Data3: 0x4519, Data4: [8]byte{0xA6, 0x63, 0x37, 0xBD, 0x56, 0x06, 0x81, 0x85}}
	FOLDERID_PrintHood              = syscall.GUID{Data1: 0x9274BD8D, Data2: 0xCFD1, Data3: 0x41C3, Data4: [8]byte{0xB3, 0x5E, 0xB1, 0x3F, 0x55, 0xA7, 0x58, 0xF4}}
	FOLDERID_Profile                = syscall.GUID{Data1: 0x5E6C858F, Data2: 0x0E22, Data3: 0x4760, Data4: [8]byte{0x9A, 0xFE, 0xEA, 0x33, 0x17, 0xB6, 0x71, 0x73}}
	FOLDERID_ProgramData            = syscall.GUID{Data1: 0x62AB5D82, Data2: 0xFDC1, Data3: 0x4DC3, Data4: [8]byte{0xA9, 0xDD, 0x07, 0x0D, 0x1D, 0x49, 0x5D, 0x97}}
	FOLDERID_ProgramFiles           = syscall.GUID{Data1: 0x905E63B6, Data2: 0xC1BF, Data3: 0x494E, Data4: [8]byte{0xB2, 0x9C, 0x65, 0xB7, 0x32, 0xD3, 0xD2, 0x1A}}
	FOLDERID_ProgramFilesX64        = syscall.GUID{Data1: 0x6D809377, Data2: 0x6AF0, Data3: 0x444B, Data4: [8]byte{0x89, 0x57, 0xA3, 0x77, 0x3F, 0x02, 0x20, 0x0E}}
	FOLDERID_ProgramFilesX86        = syscall.GUID{Data1: 0x7C5A40EF, Data2: 0xA0FB, Data3: 0x4BFC, Data4: [8]byte{0x87, 0x4A, 0xC0, 0xF2, 0xE0, 0xB9, 0xFA, 0x8E}}
	FOLDERID_ProgramFilesCommon     = syscall.GUID{Data1: 0xF7F1ED05, Data2: 0x9F6D, Data3: 0x47A2, Data4: [8]byte{0xAA, 0xAE, 0x29, 0xD3, 0x17, 0xC6, 0xF0, 0x66}}
	FOLDERID_ProgramFilesCommonX64  = syscall.GUID{Data1: 0x6365D5A7, Data2: 0x0F0D, Data3: 0x45E5, Data4: [8]byte{0x87, 0xF6, 0x0D, 0xA5, 0x6B, 0x6A, 0x4F, 0x7D}}
	FOLDERID_ProgramFilesCommonX86  = syscall.GUID{Data1: 0xDE974D24, Data2: 0xD9C6, Data3: 0x4D3E, Data4: [8]byte{0xBF, 0x91, 0xF4, 0x45, 0x51, 0x20, 0xB9, 0x17}}
	FOLDERID_Programs               = syscall.GUID{Data1: 0xA77F5D77, Data2: 0x2E2B, Data3: 0x44C3, Data4: [8]byte{0xA6, 0xA2, 0xAB, 0xA6, 0x01, 0x05, 0x4A, 0x51}}
	FOLDERID_Public                 = syscall.GUID{Data1: 0xDFDF76A2, Data2: 0xC82A, Data3: 0x4D63, Data4: [8]byte{0x90, 0x6A, 0x56, 0x44, 0xAC, 0x45, 0x73, 0x85}}
	FOLDERID_PublicDesktop          = syscall.GUID{Data1: 0xC4AA340D, Data2: 0xF20F, Data3: 0x4863, Data4: [8]byte{0xAF, 0xEF, 0xF8, 0x7E, 0xF2, 0xE6, 0xBA, 0x25}}
	FOLDERID_PublicDocuments        = syscall.GUID{Data1: 0xED4824AF, Data2: 0xDCE4, Data3: 0x45A8, Data4: [8]byte{0x81, 0xE2, 0xFC, 0x79, 0x65, 0x08, 0x36, 0x34}}
	FOLDERID_PublicDownloads        = syscall.GUID{Data1: 0x3D644C9B, Data2: 0x1FB8, Data3: 0x4F30, Data4: [8]byte{0x9B, 0x45, 0xF6, 0x70, 0x23, 0x5F, 0x79, 0xC0}}
	FOLDERID_PublicGameTasks        = syscall.GUID{Data1: 0xDEBF2536, Data2: 0xE1A8, Data3: 0x4C59, Data4: [8]byte{0xB6, 0xA2, 0x41, 0x45, 0x86, 0x47, 0x6A, 0xEA}}
	FOLDERID_PublicLibraries        = syscall.GUID{Data1: 0x48DAF80B, Data2: 0xE6CF, Data3: 0x4F4E, Data4: [8]byte{0xB8, 0x00, 0x0E, 0x69, 0xD8, 0x4E, 0xE3, 0x84}}
	FOLDERID_PublicMusic            = syscall.GUID{Data1: 0x3214FAB5, Data2: 0x9757, Data3: 0x4298, Data4: [8]byte{0xBB, 0x61, 0x92, 0xA9, 0xDE, 0xAA, 0x44, 0xFF}}
	FOLDERID_PublicPictures         = syscall.GUID{Data1: 0xB6EBFB86, Data2: 0x6907, Data3: 0x413C, Data4: [8]byte{0x9A, 0xF7, 0x4F, 0xC2, 0xAB, 0xF0, 0x7C, 0xC5}}
	FOLDERID_PublicRingtones        = syscall.GUID{Data1: 0xE555AB60, Data2: 0x153B, Data3: 0x4D17, Data4: [8]byte{0x9F, 0x04, 0xA5, 0xFE, 0x99, 0xFC, 0x15, 0xEC}}
	FOLDERID_PublicUserTiles        = syscall.GUID{Data1: 0x0482AF6C, Data2: 0x08F1, Data3: 0x4C34, Data4: [8]byte{0x8C, 0x90, 0xE1, 0x7E, 0xC9, 0x8B, 0x1E, 0x17}}
	FOLDERID_PublicVideos           = syscall.GUID{Data1: 0x2400183A, Data2: 0x6185, Data3: 0x49FB, Data4: [8]byte{0xA2, 0xD8, 0x4A, 0x39, 0x2A, 0x60, 0x2B, 0xA3}}
	FOLDERID_QuickLaunch            = syscall.GUID{Data1: 0x52A4F021, Data2: 0x7B75, Data3: 0x48A9, Data4: [8]byte{0x9F, 0x6B, 0x4B, 0x87, 0xA2, 0x10, 0xBC, 0x8F}}
	FOLDERID_Recent                 = syscall.GUID{Data1: 0xAE50C081, Data2: 0xEBD2, Data3: 0x438A, Data4: [8]byte{0x86, 0x55, 0x8A, 0x09, 0x2E, 0x34, 0x98, 0x7A}}
	FOLDERID_RecordedTVLibrary      = syscall.GUID{Data1: 0x1A6FDBA2, Data2: 0xF42D, Data3: 0x4358, Data4: [8]byte{0xA7, 0x98, 0xB7, 0x4D, 0x74, 0x59, 0x26, 0xC5}}
	FOLDERID_RecycleBinFolder       = syscall.GUID{Data1: 0xB7534046, Data2: 0x3ECB, Data3: 0x4C18, Data4: [8]byte{0xBE, 0x4E, 0x64, 0xCD, 0x4C, 0xB7, 0xD6, 0xAC}}
	FOLDERID_ResourceDir            = syscall.GUID{Data1: 0x8AD10C31, Data2: 0x2ADB, Data3: 0x4296, Data4: [8]byte{0xA8, 0xF7, 0xE4, 0x70, 0x12, 0x32, 0xC9, 0x72}}
	FOLDERID_Ringtones              = syscall.GUID{Data1: 0xC870044B, Data2: 0xF49E, Data3: 0x4126, Data4: [8]byte{0xA9, 0xC3, 0xB5, 0x2A, 0x1F, 0xF4, 0x11, 0xE8}}
	FOLDERID_RoamingAppData         = syscall.GUID{Data1: 0x3EB685DB, Data2: 0x65F9, Data3: 0x4CF6, Data4: [8]byte{0xA0, 0x3A, 0xE3, 0xEF, 0x65, 0x72, 0x9F, 0x3D}}
	FOLDERID_RoamedTileImages       = syscall.GUID{Data1: 0xAAA8D5A5, Data2: 0xF1D6, Data3: 0x4259, Data4: [8]byte{0xBA, 0xA8, 0x78, 0xE7, 0xEF, 0x60, 0x83, 0x5E}}
	FOLDERID_RoamingTiles           = syscall.GUID{Data1: 0x00BCFC5A, Data2: 0xED94, Data3: 0x4E48, Data4: [8]byte{0x96, 0xA1, 0x3F, 0x62, 0x17, 0xF2, 0x19, 0x90}}
	FOLDERID_SampleMusic            = syscall.GUID{Data1: 0xB250C668, Data2: 0xF57D, Data3: 0x4EE1, Data4: [8]byte{0xA6, 0x3C, 0x29, 0x0E, 0xE7, 0xD1, 0xAA, 0x1F}}
	FOLDERID_SamplePictures         = syscall.GUID{Data1: 0xC4900540, Data2: 0x2379, Data3: 0x4C75, Data4: [8]byte{0x84, 0x4B, 0x64, 0xE6, 0xFA, 0xF8, 0x71, 0x6B}}
	FOLDERID_SamplePlaylists        = syscall.GUID{Data1: 0x15CA69B3, Data2: 0x30EE, Data3: 0x49C1, Data4: [8]byte{0xAC, 0xE1, 0x6B, 0x5E, 0xC3, 0x72, 0xAF, 0xB5}}
	FOLDERID_SampleVideos           = syscall.GUID{Data1: 0x859EAD94, Data2: 0x2E85, Data3: 0x48AD, Data4: [8]byte{0xA7, 0x1A, 0x09, 0x69, 0xCB, 0x56, 0xA6, 0xCD}}
	FOLDERID_SavedGames             = syscall.GUID{Data1: 0x4C5C32FF, Data2: 0xBB9D, Data3: 0x43B0, Data4: [8]byte{0xB5, 0xB4, 0x2D, 0x72, 0xE5, 0x4E, 0xAA, 0xA4}}
	FOLDERID_SavedPictures          = syscall.GUID{Data1: 0x3B193882, Data2: 0xD3AD, Data3: 0x4EAB, Data4: [8]byte{0x96, 0x5A, 0x69, 0x82, 0x9D, 0x1F, 0xB5, 0x9F}}
	FOLDERID_SavedPicturesLibrary   = syscall.GUID{Data1: 0xE25B5812, Data2: 0xBE88, Data3: 0x4BD9, Data4: [8]byte{0x94, 0xB0, 0x29, 0x23, 0x34, 0x77, 0xB6, 0xC3}}
	FOLDERID_SavedSearches          = syscall.GUID{Data1: 0x7D1D3A04, Data2: 0xDEBB, Data3: 0x4115, Data4: [8]byte{0x95, 0xCF, 0x2F, 0x29, 0xDA, 0x29, 0x20, 0xDA}}
	FOLDERID_Screenshots            = syscall.GUID{Data1: 0xB7BEDE81, Data2: 0xDF94, Data3: 0x4682, Data4: [8]byte{0xA7, 0xD8, 0x57, 0xA5, 0x26, 0x20, 0xB8, 0x6F}}
	FOLDERID_SEARCH_CSC             = syscall.GUID{Data1: 0xEE32E446, Data2: 0x31CA, Data3: 0x4ABA, Data4: [8]byte{0x81, 0x4F, 0xA5, 0xEB, 0xD2, 0xFD, 0x6D, 0x5E}}
	FOLDERID_SearchHistory          = syscall.GUID{Data1: 0x0D4C3DB6, Data2: 0x03A3, Data3: 0x462F, Data4: [8]byte{0xA0, 0xE6, 0x08, 0x92, 0x4C, 0x41, 0xB5, 0xD4}}
	FOLDERID_SearchHome             = syscall.GUID{Data1: 0x190337D1, Data2: 0xB8CA, Data3: 0x4121, Data4: [8]byte{0xA6, 0x39, 0x6D, 0x47, 0x2D, 0x16, 0x97, 0x2A}}
	FOLDERID_SEARCH_MAPI            = syscall.GUID{Data1: 0x98EC0E18, Data2: 0x2098, Data3: 0x4D44, Data4: [8]byte{0x86, 0x44, 0x66, 0x97, 0x93, 0x15, 0xA2, 0x81}}
	FOLDERID_SearchTemplates        = syscall.GUID{Data1: 0x7E636BFE, Data2: 0xDFA9, Data3: 0x4D5E, Data4: [8]byte{0xB4, 0x56, 0xD7, 0xB3, 0x98, 0x51, 0xD8, 0xA9}}
	FOLDERID_SendTo                 = syscall.GUID{Data1: 0x8983036C, Data2: 0x27C0, Data3: 0x404B, Data4: [8]byte{0x8F, 0x08, 0x10, 0x2D, 0x10, 0xDC, 0xFD, 0x74}}
	FOLDERID_SidebarDefaultParts    = syscall.GUID{Data1: 0x7B396E54, Data2: 0x9EC5, Data3: 0x4300, Data4: [8]byte{0xBE, 0x0A, 0x24, 0x82, 0xEB, 0xAE, 0x1A, 0x26}}
	FOLDERID_SidebarParts           = syscall.GUID{Data1: 0xA75D362E, Data2: 0x50FC, Data3: 0x4FB7, Data4: [8]byte{0xAC, 0x2C, 0xA8, 0xBE, 0xAA, 0x31, 0x44, 0x93}}
	FOLDERID_SkyDrive               = syscall.GUID{Data1: 0xA52BBA46, Data2: 0xE9E1, Data3: 0x435F, Data4: [8]byte{0xB3, 0xD9, 0x28, 0xDA, 0xA6, 0x48, 0xC0, 0xF6}}
	FOLDERID_SkyDriveCameraRoll     = syscall.GUID{Data1: 0x767E6811, Data2: 0x49CB, Data3: 0x4273, Data4: [8]byte{0x87, 0xC2, 0x20, 0xF3, 0x55, 0xE1, 0x08, 0x5B}}
	FOLDERID_SkyDriveDocuments      = syscall.GUID{Data1: 0x24D89E24, Data2: 0x2F19, Data3: 0x4534, Data4: [8]byte{0x9D, 0xDE, 0x6A, 0x66, 0x71, 0xFB, 0xB8, 0xFE}}
	FOLDERID_SkyDrivePictures       = syscall.GUID{Data1: 0x339719B5, Data2: 0x8C47, Data3: 0x4894, Data4: [8]byte{0x94, 0xC2, 0xD8, 0xF7, 0x7A, 0xDD, 0x44, 0xA6}}
	FOLDERID_StartMenu              = syscall.GUID{Data1: 0x625B53C3, Data2: 0xAB48, Data3: 0x4EC1, Data4: [8]byte{0xBA, 0x1F, 0xA1, 0xEF, 0x41, 0x46, 0xFC, 0x19}}
	FOLDERID_Startup                = syscall.GUID{Data1: 0xB97D20BB, Data2: 0xF46A, Data3: 0x4C97, Data4: [8]byte{0xBA, 0x10, 0x5E, 0x36, 0x08, 0x43, 0x08, 0x54}}
	FOLDERID_SyncManagerFolder      = syscall.GUID{Data1: 0x43668BF8, Data2: 0xC14E, Data3: 0x49B2, Data4: [8]byte{0x97, 0xC9, 0x74, 0x77, 0x84, 0xD7, 0x84, 0xB7}}
	FOLDERID_SyncResultsFolder      = syscall.GUID{Data1: 0x289A9A43, Data2: 0xBE44, Data3: 0x4057, Data4: [8]byte{0xA4, 0x1B, 0x58, 0x7A, 0x76, 0xD7, 0xE7, 0xF9}}
	FOLDERID_SyncSetupFolder        = syscall.GUID{Data1: 0x0F214138, Data2: 0xB1D3, Data3: 0x4A90, Data4: [8]byte{0xBB, 0xA9, 0x27, 0xCB, 0xC0, 0xC5, 0x38, 0x9A}}
	FOLDERID_System                 = syscall.GUID{Data1: 0x1AC14E77, Data2: 0x02E7, Data3: 0x4E5D, Data4: [8]byte{0xB7, 0x44, 0x2E, 0xB1, 0xAE, 0x51, 0x98, 0xB7}}
	FOLDERID_SystemX86              = syscall.GUID{Data1: 0xD65231B0, Data2: 0xB2F1, Data3: 0x4857, Data4: [8]byte{0xA4, 0xCE, 0xA8, 0xE7, 0xC6, 0xEA, 0x7D, 0x27}}
	FOLDERID_Templates              = syscall.GUID{Data1: 0xA63293E8, Data2: 0x664E, Data3: 0x48DB, Data4: [8]byte{0xA0, 0x79, 0xDF, 0x75, 0x9E, 0x05, 0x09, 0xF7}}
	FOLDERID_UserPinned             = syscall.GUID{Data1: 0x9E3995AB, Data2: 0x1F9C, Data3: 0x4F13, Data4: [8]byte{0xB8, 0x27, 0x48, 0xB2, 0x4B, 0x6C, 0x71, 0x74}}
	FOLDERID_UserProfiles           = syscall.GUID{Data1: 0x0762D272, Data2: 0xC50A, Data3: 0x4BB0, Data4: [8]byte{0xA3, 0x82, 0x69, 0x7D, 0xCD, 0x72, 0x9B, 0x80}}
	FOLDERID_UserProgramFiles       = syscall.GUID{Data1: 0x5CD7AEE2, Data2: 0x2219, Data3: 0x4A67, Data4: [8]byte{0xB8, 0x5D, 0x6C, 0x9C, 0xE1, 0x56, 0x60, 0xCB}}
	FOLDERID_UserProgramFilesCommon = syscall.GUID{Data1: 0xBCBD3057, Data2: 0xCA5C, Data3: 0x4622, Data4: [8]byte{0xB4, 0x2D, 0xBC, 0x56, 0xDB, 0x0A, 0xE5, 0x16}}
	FOLDERID_UsersFiles             = syscall.GUID{Data1: 0xF3CE0F7C, Data2: 0x4901, Data3: 0x4ACC, Data4: [8]byte{0x86, 0x48, 0xD5, 0xD4, 0x4B, 0x04, 0xEF, 0x8F}}
	FOLDERID_UsersLibraries         = syscall.GUID{Data1: 0xA302545D, Data2: 0xDEFF, Data3: 0x464B, Data4: [8]byte{0xAB, 0xE8, 0x61, 0xC8, 0x64, 0x8D, 0x93, 0x9B}}
	FOLDERID_Videos                 = syscall.GUID{Data1: 0x18989B1D, Data2: 0x99B5, Data3: 0x455B, Data4: [8]byte{0x84, 0x1C, 0xAB, 0x7C, 0x74, 0xE4, 0xDD, 0xFC}}
	FOLDERID_VideosLibrary          = syscall.GUID{Data1: 0x491E922F, Data2: 0x5643, Data3: 0x4AF4, Data4: [8]byte{0xA7, 0xEB, 0x4E, 0x7A, 0x13, 0x8D, 0x81, 0x74}}
	FOLDERID_Windows                = syscall.GUID{Data1: 0xF38BF404, Data2: 0x1D43, Data3: 0x42F2, Data4: [8]byte{0x93, 0x05, 0x67, 0xDE, 0x0B, 0x28, 0xFC, 0x23}}
)

type Msv1_0_LogonSubmitType uint32
type Msv1_0_ProtocolMessageType uint32

var (
	KF_FLAG_DEFAULT                     uint32 = 0x00000000
	KF_FLAG_SIMPLE_IDLIST               uint32 = 0x00000100
	KF_FLAG_NOT_PARENT_RELATIVE         uint32 = 0x00000200
	KF_FLAG_DEFAULT_PATH                uint32 = 0x00000400
	KF_FLAG_INIT                        uint32 = 0x00000800
	KF_FLAG_NO_ALIAS                    uint32 = 0x00001000
	KF_FLAG_DONT_UNEXPAND               uint32 = 0x00002000
	KF_FLAG_DONT_VERIFY                 uint32 = 0x00004000
	KF_FLAG_CREATE                      uint32 = 0x00008000
	KF_FLAG_NO_APPCONTAINER_REDIRECTION uint32 = 0x00010000
	KF_FLAG_ALIAS_ONLY                  uint32 = 0x80000000

	MICROSOFT_KERBEROS_NAME_A LSAString = LSAStringMustCompile("Kerberos")
	MSV1_0_PACKAGE_NAME       LSAString = LSAStringMustCompile("MICROSOFT_AUTHENTICATION_PACKAGE_V1_0")
	NEGOSSP_NAME_A            LSAString = LSAStringMustCompile("Negotiate")

	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378764(v=vs.85).aspx
	// MSV1_0_LOGON_SUBMIT_TYPE enumeration
	MsV1_0InteractiveLogon       Msv1_0_LogonSubmitType = 2
	MsV1_0Lm20Logon              Msv1_0_LogonSubmitType = 3
	MsV1_0NetworkLogon           Msv1_0_LogonSubmitType = 4
	MsV1_0SubAuthLogon           Msv1_0_LogonSubmitType = 5
	MsV1_0WorkstationUnlockLogon Msv1_0_LogonSubmitType = 7
	MsV1_0S4ULogon               Msv1_0_LogonSubmitType = 12
	MsV1_0VirtualLogon           Msv1_0_LogonSubmitType = 82
	MsV1_0NoElevationLogon       Msv1_0_LogonSubmitType = 82

	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378766(v=vs.85).aspx
	// MSV1_0_PROTOCOL_MESSAGE_TYPE
	MsV1_0Lm20ChallengeRequest     Msv1_0_ProtocolMessageType = 0
	MsV1_0Lm20GetChallengeResponse Msv1_0_ProtocolMessageType = 1
	MsV1_0EnumerateUsers           Msv1_0_ProtocolMessageType = 2
	MsV1_0GetUserInfo              Msv1_0_ProtocolMessageType = 3
	MsV1_0ReLogonUsers             Msv1_0_ProtocolMessageType = 4
	MsV1_0ChangePassword           Msv1_0_ProtocolMessageType = 5
	MsV1_0ChangeCachedPassword     Msv1_0_ProtocolMessageType = 6
	MsV1_0GenericPassthrough       Msv1_0_ProtocolMessageType = 7
	MsV1_0CacheLogon               Msv1_0_ProtocolMessageType = 8
	MsV1_0SubAuth                  Msv1_0_ProtocolMessageType = 9
	MsV1_0DeriveCredential         Msv1_0_ProtocolMessageType = 10
	MsV1_0CacheLookup              Msv1_0_ProtocolMessageType = 11
	MsV1_0SetProcessOption         Msv1_0_ProtocolMessageType = 12
	MsV1_0ConfigLocalAliases       Msv1_0_ProtocolMessageType = 13
	MsV1_0ClearCachedCredentials   Msv1_0_ProtocolMessageType = 14
	MsV1_0LookupToken              Msv1_0_ProtocolMessageType = 15
	MsV1_0ValidateAuth             Msv1_0_ProtocolMessageType = 16
	MsV1_0CacheLookupEx            Msv1_0_ProtocolMessageType = 17
	MsV1_0GetCredentialKey         Msv1_0_ProtocolMessageType = 18
	MsV1_0SetThreadOption          Msv1_0_ProtocolMessageType = 19
)

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686219(v=vs.85).aspx
	ABOVE_NORMAL_PRIORITY_CLASS   = 0x00008000
	BELOW_NORMAL_PRIORITY_CLASS   = 0x00004000
	HIGH_PRIORITY_CLASS           = 0x00000080
	IDLE_PRIORITY_CLASS           = 0x00000040
	NORMAL_PRIORITY_CLASS         = 0x00000020
	PROCESS_MODE_BACKGROUND_BEGIN = 0x00100000
	PROCESS_MODE_BACKGROUND_END   = 0x00200000
	REALTIME_PRIORITY_CLASS       = 0x00000100
)

func LSAStringMustCompile(s string) LSAString {
	l, err := LSAStringFromString(s)
	if err != nil {
		panic(err)
	}
	return l
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686347(v=vs.85).aspx
func SwitchDesktop(
	hDesktop Hdesk, // HDESK
) (err error) {
	r1, _, e1 := procSwitchDesktop.Call(
		uintptr(hDesktop),
	)
	if r1 == 0 {
		err = os.NewSyscallError("SwitchDesktop", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682024(v=vs.85).aspx
func CloseDesktop(
	hDesktop Hdesk, // HDESK
) (err error) {
	r1, _, e1 := procCloseDesktop.Call(
		uintptr(hDesktop),
	)
	if r1 == 0 {
		err = os.NewSyscallError("CloseDesktop", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686219(v=vs.85).aspx
func SetPriorityClass(
	hProcess syscall.Handle, // HANDLE
	dwPriorityClass uint32, // DWORD
) (err error) {
	r1, _, e1 := procSetPriorityClass.Call(
		uintptr(hProcess),
		uintptr(dwPriorityClass),
	)
	if r1 == 0 {
		err = os.NewSyscallError("SetPriorityClass", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762270(v=vs.85).aspx
func CreateEnvironmentBlock(
	lpEnvironment *uintptr, // LPVOID*
	hToken syscall.Handle, // HANDLE
	bInherit bool, // BOOL
) (err error) {
	inherit := uint32(0)
	if bInherit {
		inherit = 1
	}
	r1, _, e1 := procCreateEnvironmentBlock.Call(
		uintptr(unsafe.Pointer(lpEnvironment)),
		uintptr(hToken),
		uintptr(inherit),
	)
	if r1 == 0 {
		err = os.NewSyscallError("CreateEnvironmentBlock", e1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762274(v=vs.85).aspx
func DestroyEnvironmentBlock(
	lpEnvironment uintptr, // LPVOID - beware - unlike LPVOID* in CreateEnvironmentBlock!
) (err error) {
	r1, _, e1 := procDestroyEnvironmentBlock.Call(
		lpEnvironment,
	)
	if r1 == 0 {
		err = os.NewSyscallError("DestroyEnvironmentBlock", e1)
	}
	return
}

// CreateEnvironment returns an environment block, suitable for use with the
// CreateProcessAsUser system call. The default environment variables of hUser
// are overlayed with values in env.
func CreateEnvironment(env *[]string, hUser syscall.Handle) (envBlock *uint16, err error) {
	var logonEnv uintptr
	err = CreateEnvironmentBlock(&logonEnv, hUser, false)
	if err != nil {
		return
	}
	defer DestroyEnvironmentBlock(logonEnv)
	var varStartOffset uint
	envList := &[]string{}
	for {
		envVar := syscall.UTF16ToString((*[1 << 15]uint16)(unsafe.Pointer(logonEnv + uintptr(varStartOffset)))[:])
		if envVar == "" {
			break
		}
		*envList = append(*envList, envVar)
		// in UTF16, each rune takes two bytes, as does the trailing uint16(0)
		varStartOffset += uint(2 * (utf8.RuneCountInString(envVar) + 1))
	}
	env, err = MergeEnvLists(envList, env)
	if err != nil {
		return
	}
	return ListToEnvironmentBlock(env), nil
}

type envSetting struct {
	name  string
	value string
}

func MergeEnvLists(envLists ...*[]string) (*[]string, error) {
	mergedEnvMap := map[string]envSetting{}
	for _, envList := range envLists {
		if envList == nil {
			continue
		}
		for _, env := range *envList {
			if utf8.RuneCountInString(env) > 32767 {
				return nil, fmt.Errorf("Env setting is more than 32767 runes: %v", env)
			}
			spl := strings.SplitN(env, "=", 2)
			if len(spl) != 2 {
				return nil, fmt.Errorf("Could not interpret string %q as `key=value`", env)
			}
			newVarName := spl[0]
			newVarValue := spl[1]
			// if env var already exists, use case of existing name, to simulate behaviour of
			// setting an existing env var with a different case
			// e.g.
			//  set aVar=3
			//  set AVAR=4
			// results in
			//  aVar=4
			canonicalVarName := strings.ToLower(newVarName)
			if existingVarName := mergedEnvMap[canonicalVarName].name; existingVarName != "" {
				newVarName = existingVarName
			}
			mergedEnvMap[canonicalVarName] = envSetting{
				name:  newVarName,
				value: newVarValue,
			}
		}
	}
	canonicalVarNames := make([]string, len(mergedEnvMap))
	i := 0
	for k := range mergedEnvMap {
		canonicalVarNames[i] = k
		i++
	}
	// All strings in the environment block must be sorted alphabetically by
	// name. The sort is case-insensitive, Unicode order, without regard to
	// locale.
	//
	// See https://msdn.microsoft.com/en-us/library/windows/desktop/ms682009(v=vs.85).aspx
	sort.Strings(canonicalVarNames)
	// Finally piece back together into an environment block
	mergedEnv := make([]string, len(mergedEnvMap))
	i = 0
	for _, canonicalVarName := range canonicalVarNames {
		mergedEnv[i] = mergedEnvMap[canonicalVarName].name + "=" + mergedEnvMap[canonicalVarName].value
		i++
	}
	return &mergedEnv, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762188(v=vs.85).aspx
func SHGetKnownFolderPath(rfid *syscall.GUID, dwFlags uint32, hToken syscall.Handle, pszPath *uintptr) (err error) {
	r0, _, _ := procSHGetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762249(v=vs.85).aspx
func SHSetKnownFolderPath(
	rfid *syscall.GUID, // REFKNOWNFOLDERID
	dwFlags uint32, // DWORD
	hToken syscall.Handle, // HANDLE
	pszPath *uint16, // PCWSTR
) (err error) {
	r1, _, _ := procSHSetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(rfid)),
		uintptr(dwFlags),
		uintptr(hToken),
		uintptr(unsafe.Pointer(pszPath)),
	)
	if r1 != 0 {
		err = syscall.Errno(r1)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms680722(v=vs.85).aspx
// Note: the system call returns no value, so we can't check for an error
func CoTaskMemFree(pv uintptr) {
	procCoTaskMemFree.Call(uintptr(pv))
}

func CloseHandle(handle syscall.Handle) (err error) {
	syscall.CloseHandle(handle)
	r1, _, e1 := procCloseHandle.Call(
		uintptr(handle),
	)
	if r1 == 0 {
		if e1 != syscall.Errno(0) {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		log.Printf("Error when closing handle: %v", err)
	}
	return
}

func GetFolder(hUser syscall.Handle, folder *syscall.GUID, dwFlags uint32) (value string, err error) {
	var path uintptr
	err = SHGetKnownFolderPath(folder, dwFlags, hUser, &path)
	if err != nil {
		return
	}
	// CoTaskMemFree system call has no return value, so can't check for error
	defer CoTaskMemFree(path)
	value = syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(path))[:])
	return
}

func SetFolder(hUser syscall.Handle, folder *syscall.GUID, value string) (err error) {
	var s *uint16
	s, err = syscall.UTF16PtrFromString(value)
	if err != nil {
		return
	}
	return SHSetKnownFolderPath(folder, 0, hUser, s)
}

func SetAndCreateFolder(hUser syscall.Handle, folder *syscall.GUID, value string) (err error) {
	err = SetFolder(hUser, folder, value)
	if err != nil {
		return
	}
	_, err = GetFolder(hUser, folder, KF_FLAG_CREATE)
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378265(v=vs.85).aspx
func LsaConnectUntrusted(
	lsaHandle *syscall.Handle, // PHANDLE
) (err error) {
	r, _, e := procLsaConnectUntrusted.Call(
		uintptr(unsafe.Pointer(lsaHandle)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaConnectUntrusted sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378522(v=vs.85).aspx
type LSAString struct {
	Length        uint16 // USHORT
	MaximumLength uint16 // USHORT
	Buffer        *uint8 // PCHAR
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378522(v=vs.85).aspx
func LSAStringFromString(s string) (LSAString, error) {
	ansi, err := charmap.Windows1252.NewEncoder().Bytes([]byte(s))
	if err != nil {
		return LSAString{}, err
	}
	return LSAString{
		Length:        uint16(len(ansi)),
		MaximumLength: uint16(len(ansi)),
		Buffer:        &ansi[0],
	}, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378297(v=vs.85).aspx
func LsaLookupAuthenticationPackage(
	lsaHandle syscall.Handle, // HANDLE
	packageName *LSAString, // PLSA_STRING
	authenticationPackage *uint32, // PULONG
) (err error) {
	r, _, e := procLsaLookupAuthenticationPackage.Call(
		uintptr(lsaHandle),
		uintptr(unsafe.Pointer(packageName)),
		uintptr(unsafe.Pointer(authenticationPackage)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaLookupAuthenticationPackage sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379261(v=vs.85).aspx
type LUID struct {
	LowPart  uint32 // DWORD
	HighPart int32  // LONG
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa375260(v=vs.85).aspx
func AllocateLocallyUniqueId(
	luid *LUID, // PLUID
) (err error) {
	r, _, _ := procAllocateLocallyUniqueId.Call(
		uintptr(unsafe.Pointer(luid)),
	)
	if r == 0 {
		err = syscall.Errno(r)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379631(v=vs.85).aspx
type TokenSource struct {
	SourceName       [8]byte // CHAR[TOKEN_SOURCE_LENGTH]
	SourceIdentifier LUID    // LUID
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379594(v=vs.85).aspx
type PSID uintptr

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379595(v=vs.85).aspx
type SidAndAttributes struct {
	Sid        PSID   // PSID
	Attributes uint32 // DWORD
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379624(v=vs.85).aspx
type TokenGroups struct {
	GroupCount uint32           // DWORD
	Groups     SidAndAttributes // SID_AND_ATTRIBUTES[ANYSIZE_ARRAY]
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa379363(v=vs.85).aspx
type QuotaLimits struct {
	PagedPoolLimit        uintptr // SIZE_T
	NonPagedPoolLimit     uintptr // SIZE_T
	MinimumWorkingSetSize uintptr // SIZE_T
	MaximumWorkingSetSize uintptr // SIZE_T
	PagefileLimit         uintptr // SIZE_T
	TimeLimit             uint64  // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
}

// https://msdn.microsoft.com/en-us/library/windows/hardware/ff565436(v=vs.85).aspx
type NtStatus uint32

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa380129(v=vs.85).aspx
// Below, "uint32" is determined by experimenting with C++, since
// https://msdn.microsoft.com/en-us/library/2dzy4k6e.aspx simply says "The
// underlying type of the enumerators; all enumerators have the same underlying
// type. May be any integral type." - thank you Microsoft and C++.
// By writing a C++ program that displays the memory, I determined it is 4 bytes.
type SecurityLogonType uint32

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378093(v=vs.85).aspx
type KerbLogonSubmitType uint

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378096(v=vs.85).aspx
type KerbProfileBufferType uint

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378079(v=vs.85).aspx
type KerbInteractiveLogon struct {
	MessageType     KerbLogonSubmitType
	LogonDomainName ntr.LSAUnicodeString
	UserName        ntr.LSAUnicodeString
	Password        ntr.LSAUnicodeString
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378147(v=vs.85).aspx
type KerbTicketProfile struct {
	Profile    KerbInteractiveProfile
	SessionKey KerbCryptoKey
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378083(v=vs.85).aspx
type KerbInteractiveProfile struct {
	MessageType        KerbProfileBufferType
	LogonCount         uint16               // USHORT
	BadPasswordCount   uint16               // USHORT
	LogonTime          uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	LogoffTime         uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	KickOffTime        uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	PasswordLastSet    uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	PasswordCanChange  uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	PasswordMustChange uint64               // LARGE_INTEGER -- is actually a union - might cause trouble on win32, *sigh*
	LogonScript        ntr.LSAUnicodeString // UNICODE_STRING
	HomeDirectory      ntr.LSAUnicodeString // UNICODE_STRING
	FullName           ntr.LSAUnicodeString // UNICODE_STRING
	ProfilePath        ntr.LSAUnicodeString // UNICODE_STRING
	HomeDirectoryDrive ntr.LSAUnicodeString // UNICODE_STRING
	LogonServer        ntr.LSAUnicodeString // UNICODE_STRING
	UserFlags          uint32               // ULONG
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378058(v=vs.85).aspx
type KerbCryptoKey struct {
	KeyType int32  // LONG
	Length  uint32 // ULONG
	Value   *uint8 // PUCHAR
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378292(v=vs.85).aspx
func LsaLogonUser(
	lsaHandle syscall.Handle, // HANDLE
	originName *LSAString, // PLSA_STRING
	logonType SecurityLogonType, // SECURITY_LOGON_TYPE
	authenticationPackage uint32, // ULONG
	authenticationInformation *byte, // PVOID
	authenticationInformationLength uint32, // ULONG
	localGroups *TokenGroups, // PTOKEN_GROUPS
	sourceContext *TokenSource, // PTOKEN_SOURCE
	profileBuffer *uintptr, // PVOID*
	profileBufferLength *uint32, // PULONG
	logonId *LUID, // PLUID
	token *syscall.Handle, // PHANDLE
	quotas *QuotaLimits, // PQUOTA_LIMITS
	subStatus *NtStatus, // PNTSTATUS
) (err error) {
	r, _, e := procLsaLogonUser.Call(
		uintptr(lsaHandle),
		uintptr(unsafe.Pointer(originName)),
		uintptr(logonType),
		uintptr(authenticationPackage),
		uintptr(unsafe.Pointer(authenticationInformation)),
		uintptr(authenticationInformationLength),
		uintptr(unsafe.Pointer(localGroups)),
		uintptr(unsafe.Pointer(sourceContext)),
		uintptr(unsafe.Pointer(profileBuffer)),
		uintptr(unsafe.Pointer(profileBufferLength)),
		uintptr(unsafe.Pointer(logonId)),
		uintptr(unsafe.Pointer(token)),
		uintptr(unsafe.Pointer(quotas)),
		uintptr(unsafe.Pointer(subStatus)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaLookupAuthenticationPackage sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

// TODO: nice-to-have but not necessary: would make troubleshooting easier
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms721800(v=vs.85).aspx
func LsaNtStatusToWinError() {}

// InteractiveDomainLogon authenticates a security principal's logon data
// and returns a new interactive logon session.
//
// The originName parameter should specify meaningful information. For example,
// it might contain "TTY1" to indicate terminal one or "NTLM - remote node
// JAZZ" to indicate a network logon that uses NTLM through a remote node
// called "JAZZ".
//
// The sourceName parameter specifies an 8-byte ASCII character string used to
// identify the source of an access token. This is used to distinguish between
// such sources as Session Manager, LAN Manager, and RPC Server. A string,
// rather than a constant, is used to identify the source so users and
// developers can make extensions to the system, such as by adding other
// networks, that act as the source of access tokens.
func InteractiveDomainLogon(logonDomain, username, password, originName string, sourceName [8]byte) (profileBuffer uintptr, profileBufferLength uint32, logonId LUID, token syscall.Handle, quotas QuotaLimits, subStatus NtStatus, err error) {

	// Before making any syscalls, first validate inputs...

	var lsaDomain, lsaUsername, lsaPassword ntr.LSAUnicodeString

	lsaDomain, err = ntr.LSAUnicodeStringFromString(logonDomain)
	if err != nil {
		return
	}
	lsaUsername, err = ntr.LSAUnicodeStringFromString(username)
	if err != nil {
		return
	}
	lsaPassword, err = ntr.LSAUnicodeStringFromString(password)
	if err != nil {
		return
	}
	var oName LSAString
	oName, err = LSAStringFromString(originName)
	if err != nil {
		return
	}

	// Validation complete.

	h := syscall.Handle(0)
	err = LsaConnectUntrusted(&h)
	if err != nil {
		return
	}

	authPackage := uint32(0)
	err = LsaLookupAuthenticationPackage(h, &MICROSOFT_KERBEROS_NAME_A, &authPackage)
	if err != nil {
		return
	}

	l := LUID{}
	err = AllocateLocallyUniqueId(&l)
	if err != nil {
		return
	}

	authenticationInformation := KerbInteractiveLogon{
		MessageType:     2, // KerbInteractiveLogon
		LogonDomainName: lsaDomain,
		UserName:        lsaUsername,
		Password:        lsaPassword,
	}

	// Although the MSDN docs don't state this, the authenticationInformation
	// buffer has to be a contiguous block that includes the KerbInteractiveLogon
	// struct and the three LSA unicode string buffers.

	// This requires some delicate surgery.

	// 1) Copy the three LSA unicode string buffers into their own []byte
	domainNameCopy := (*(*[1<<31 - 1]byte)(unsafe.Pointer(lsaDomain.Buffer)))[:lsaDomain.MaximumLength]
	userNameCopy := (*(*[1<<31 - 1]byte)(unsafe.Pointer(lsaUsername.Buffer)))[:lsaUsername.MaximumLength]
	passwordCopy := (*(*[1<<31 - 1]byte)(unsafe.Pointer(lsaPassword.Buffer)))[:lsaPassword.MaximumLength]

	// 2) Calculate the number of bytes of contiguous memory required to fit the
	// struct and the three LSA unicode string buffers
	authenticationInformationLength := uint32(unsafe.Sizeof(authenticationInformation)) + uint32(len(domainNameCopy)+len(userNameCopy)+len(passwordCopy))

	// 3) Allocate a new []byte to store the entire block
	authInfoBuffer := make([]byte, authenticationInformationLength)

	// 4) Calculate the physical address inside this new block where the LSA
	// buffers will sit, and update the existing pointer references to the new
	// locations (even though the buffers do not exist there yet)
	authenticationInformation.LogonDomainName.Buffer = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(&authInfoBuffer[0])) + unsafe.Sizeof(KerbInteractiveLogon{})))
	authenticationInformation.UserName.Buffer = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(authenticationInformation.LogonDomainName.Buffer)) + uintptr(authenticationInformation.LogonDomainName.MaximumLength)))
	authenticationInformation.Password.Buffer = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(authenticationInformation.UserName.Buffer)) + uintptr(authenticationInformation.UserName.MaximumLength)))

	// 5) Copy the updated authenticationInformation (with correct absolute pointer
	// references to LSA buffers) into a new []byte
	authInfoCopy := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&authenticationInformation)))[:unsafe.Sizeof(authenticationInformation)]

	// 6) Append the three LSA buffers to it
	authInfoCopy = append(authInfoCopy, domainNameCopy...)
	authInfoCopy = append(authInfoCopy, userNameCopy...)
	authInfoCopy = append(authInfoCopy, passwordCopy...)

	// 7) Copy the entire block over to the target block we allocated in step 3
	copy(authInfoBuffer, authInfoCopy)

	err = LsaLogonUser(
		h,
		&oName,
		2, // Interactive
		authPackage,
		&authInfoBuffer[0],
		authenticationInformationLength,
		nil,
		&TokenSource{
			SourceName:       sourceName,
			SourceIdentifier: l,
		},
		&profileBuffer,
		&profileBufferLength,
		&logonId,
		&token,
		&quotas,
		&subStatus,
	)
	// Not sure if err != nil can leave corrupt data, so to be safe, return newly allocated zero values
	if err != nil {
		return 0, 0, LUID{}, 0, QuotaLimits{}, 0, err
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378261(v=vs.85).aspx
func LsaCallAuthenticationPackage(
	lsaHandle syscall.Handle, // HANDLE
	authenticationPackage uint32, // ULONG
	protocolSubmitBuffer *byte, // PVOID
	submitBufferLength uint32, // ULONG
	protocolReturnBuffer *uintptr, // *PVOID
	returnBufferLength **uint32, // *PULONG
	protocolStatus *NtStatus, // PNTSTATUS
) (err error) {
	r, _, e := procLsaCallAuthenticationPackage.Call(
		uintptr(lsaHandle),
		uintptr(authenticationPackage),
		uintptr(unsafe.Pointer(protocolSubmitBuffer)),
		uintptr(submitBufferLength),
		uintptr(unsafe.Pointer(protocolReturnBuffer)),
		uintptr(unsafe.Pointer(returnBufferLength)),
		uintptr(unsafe.Pointer(protocolStatus)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaCallAuthenticationPackage sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378279(v=vs.85).aspx
func LsaFreeReturnBuffer(
	buffer *byte, // PVOID
) (err error) {
	r, _, e := procLsaFreeReturnBuffer.Call(
		uintptr(unsafe.Pointer(buffer)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaFreeReturnBuffer sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

type LSAOperationalMode uint32 // ULONG

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378318(v=vs.85).aspx
func LsaRegisterLogonProcess(
	logonProcessName *LSAString, // PLSA_STRING
	lsaHandle *syscall.Handle, // PHANDLE
	securityMode *LSAOperationalMode, // PLSA_OPERATIONAL_MODE
) (err error) {
	r, _, e := procLsaRegisterLogonProcess.Call(
		uintptr(unsafe.Pointer(logonProcessName)),
		uintptr(unsafe.Pointer(lsaHandle)),
		uintptr(unsafe.Pointer(securityMode)),
	)
	if r != 0 {
		err = fmt.Errorf("Got error from LsaRegisterLogonProcess sys call: 0x%X, see https://msdn.microsoft.com/en-us/library/cc704588.aspx for details: %v", r, e)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378767(v=vs.85).aspx
// typedef struct _MSV1_0_SUBAUTH_LOGON{
//     MSV1_0_LOGON_SUBMIT_TYPE MessageType;
//     UNICODE_STRING LogonDomainName;
//     UNICODE_STRING UserName;
//     UNICODE_STRING Workstation;
//     UCHAR ChallengeToClient[MSV1_0_CHALLENGE_LENGTH];
//     STRING AuthenticationInfo1;
//     STRING AuthenticationInfo2;
//     ULONG ParameterControl;
//     ULONG SubAuthPackageId;
// } MSV1_0_SUBAUTH_LOGON, * PMSV1_0_SUBAUTH_LOGON;
type Msv1_0_SubAuthLogon struct {
	MessageType         Msv1_0_LogonSubmitType
	LogonDomainName     ntr.LSAUnicodeString
	UserName            ntr.LSAUnicodeString
	Workstation         ntr.LSAUnicodeString
	CallengeToClient    [8]byte
	AuthenticationInfo1 LSAString
	AuthenticationInfo2 LSAString
	ParameterControl    uint32
	SubAuthPackageId    uint32
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378762(v=vs.85).aspx
// typedef struct _MSV1_0_LM20_LOGON {
//   MSV1_0_LOGON_SUBMIT_TYPE MessageType;
//   UNICODE_STRING           LogonDomainName;
//   UNICODE_STRING           UserName;
//   UNICODE_STRING           Workstation;
//   UCHAR                    ChallengeToClient[MSV1_0_CHALLENGE_LENGTH];
//   STRING                   CaseSensitiveChallengeResponse;
//   STRING                   CaseInsensitiveChallengeResponse;
//   ULONG                    ParameterControl;
// } MSV1_0_LM20_LOGON, *PMSV1_0_LM20_LOGON;
type Msv1_0_Lm20Logon struct {
	MessageType                      Msv1_0_LogonSubmitType
	LogonDomainName                  ntr.LSAUnicodeString
	UserName                         ntr.LSAUnicodeString
	Workstation                      ntr.LSAUnicodeString
	CallengeToClient                 [8]byte
	CaseSensitiveChallengeResponse   LSAString
	CaseInsensitiveChallengeResponse LSAString
	ParameterControl                 uint32
}

// TODO: https://msdn.microsoft.com/en-us/library/windows/desktop/aa378269(v=vs.85).aspx
func LsaDeregisterLogonProcess() {
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa378760(v=vs.85).aspx
// typedef struct _MSV1_0_INTERACTIVE_LOGON {
//   MSV1_0_LOGON_SUBMIT_TYPE MessageType;
//   UNICODE_STRING           LogonDomainName;
//   UNICODE_STRING           UserName;
//   UNICODE_STRING           Password;
// } MSV1_0_INTERACTIVE_LOGON;
type Msv1_0_InteractiveLogon struct {
	MessageType     Msv1_0_LogonSubmitType
	LogonDomainName ntr.LSAUnicodeString
	UserName        ntr.LSAUnicodeString
	Password        ntr.LSAUnicodeString
}
