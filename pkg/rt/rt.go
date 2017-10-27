package rt

type GwkRT struct {
}

var rt = &GwkRT{}

func GetRTInstance() *GwkRT {
	return rt
}

func (rt *GwkRT) init() {

}

func (rt *GwkRT) GetViewPort() (int, int) {
	width := 100
	height := 200

	return width, height
}
