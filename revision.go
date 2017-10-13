package rx

//jig:template RxRevision
//jig:revision

// RxRevision_d4a045e is the revision identifier that changes whenever the rx
// library gets updated for every (even minor) change. This is because the
// library is used by copying fragments of the source of the library instead
// of using a compiled version. When compiling the program that uses an older
// revision of rx, the build will fail and the user will have to run `jig -r`
func RxRevision_d4a045e() string { return "v0.9" }
