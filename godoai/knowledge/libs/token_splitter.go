package libs

// SplitText splits a text into multiple text.
func SplitTokenText(text string, contextLength int) ([]string, error) {

	texts := splitText(text, contextLength)

	return texts, nil
}
func splitText(text string, ChunkSize int) []string {
	splits := make([]string, 0)
	Texts := SplitText(text, ChunkSize)
	// for _, txt := range Texts {
	// 	splits = append(splits, txt)
	// }
	splits = append(splits, Texts...)
	return splits

}

/*
func spText(text string, tk *tokenizers.Tokenizer) string {
	//splits := make([]string, 0)
	_, tokens := tk.Encode(text, false)
	//fmt.Printf("tokens==%v", tokens)

	if len(tokens) == 0 {
		return "" // 如果没有tokens，直接返回空切片
	}
	//re := regexp.MustCompile(`\[[Uu][Nn][Kk]\]|\[[Pp][Aa][Dd]\]`)
	//str := re.ReplaceAllString(strings.Join(tokens, ""), "")
	var str string
	for _, token := range tokens {
		if token == "[UNK]" || token == "[PAD]" {
			continue
		}
		str += token
	}
	str = strings.ReplaceAll(str, "##", "")
	return str
}
*/
