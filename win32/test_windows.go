package win32

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/taskcluster/ntr"
)

func main() {
	h := syscall.Handle(0)
	err := LsaConnectUntrusted(&h)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Handle: %v", h)
	authPackage := uint32(0)
	err = LsaLookupAuthenticationPackage(h, &MICROSOFT_KERBEROS_NAME_A, &authPackage)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Auth package: %v", authPackage)

	l := LUID{}
	err = AllocateLocallyUniqueId(&l)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("LUID: %#v", l)

	authenticationInformation := KerbInteractiveLogon{
		MessageType:     2, // KerbInteractiveLogon
		LogonDomainName: ntr.LSAUnicodeStringMustCompile("."),
		UserName:        ntr.LSAUnicodeStringMustCompile("testuser"),
		Password:        ntr.LSAUnicodeStringMustCompile("testpassword"),
	}

	originName := LSAStringMustCompile("TestAppFoo")

	sourceName := [8]byte{}
	for i, c := range []byte("foobarxx") {
		sourceName[i] = c
	}
	for i, c := range sourceName {
		log.Printf("%v: #%v", i, c)
	}

	logonId := LUID{}
	token := syscall.Handle(0)
	quotas := QuotaLimits{}
	subStatus := NtStatus(0)

	profileBufferLength := uint32(0)

	var profileBuffer uintptr

	LsaLogonUser(
		h,
		&originName,
		2, // Interactive
		authPackage,
		&authenticationInformation,
		uint32(unsafe.Sizeof(authenticationInformation)),
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
}
