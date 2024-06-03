package commander

func Tokenize(line string) []string {
	tokens := []string{}
	curTok := []byte{}
	var quote byte = 0
	in := false
	for i := 0; i < len(line); i++ {
		if !in && line[i] == ' ' || line[i] == '\t' {
			continue
		}

		if !in && line[i] == 39 || line[i] == '"' {
			quote = line[i]
			in = true
			continue
		}

		if quote != 0 && line[i] == quote {
			quote = 0
			continue
		}

		if in && quote == 0 && line[i] == ' ' || line[i] == '\t' {
			in = false
			tokens = append(tokens, string(curTok))
			curTok = []byte{}
			continue
		}

		in = true
		curTok = append(curTok, line[i])
	}

	if len(curTok) > 0 {
		tokens = append(tokens, string(curTok))
	}

	return tokens
}
