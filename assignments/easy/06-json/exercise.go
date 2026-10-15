package lab

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func DecodeProfile(data []byte) (Profile, error) { panic("TODO: DecodeProfile") }
