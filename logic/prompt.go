package logic

type PromptLogic struct {
	Prompt string
}

func (pl *PromptLogic) GenerateImage() (string, error) {
	// TODO generate image
	imageUrl := "http://hoge.png"
	return imageUrl, nil
}
