package commander

var (
	FLOW_CONTROL_UNSPECIFIED byte = ' '
	FLOW_CONTROL_PIPE        byte = '|'
	FLOW_CONTROL_REDIRECT    byte = '>'
)

var ControlCharacters = map[byte]struct{}{
	FLOW_CONTROL_PIPE:     {},
	FLOW_CONTROL_REDIRECT: {},
}

type TokenGroup struct {
	Tokens      []string
	FlowControl byte
}

func Tokenize(line string) []*TokenGroup {
	allTokens := []*TokenGroup{}
	tokenGroup := &TokenGroup{
		Tokens:      []string{},
		FlowControl: FLOW_CONTROL_UNSPECIFIED,
	}
	curTok := []byte{}
	var quote byte = 0
	in := false
	for i := 0; i < len(line); i++ {
		if quote == 0 {
			if _, has := ControlCharacters[line[i]]; has {
				allTokens = append(allTokens, tokenGroup)
				tokenGroup = &TokenGroup{
					Tokens:      []string{},
					FlowControl: line[i],
				}
				continue
			}
		}

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
			tokenGroup.Tokens = append(tokenGroup.Tokens, string(curTok))
			curTok = []byte{}
			continue
		}

		in = true
		curTok = append(curTok, line[i])
	}

	if len(curTok) > 0 {
		tokenGroup.Tokens = append(tokenGroup.Tokens, string(curTok))
	}

	if len(tokenGroup.Tokens) > 0 {
		allTokens = append(allTokens, tokenGroup)
	}

	return allTokens
}
