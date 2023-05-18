package wordle

import (
	words "golang-addon/week-1/golang-wordle/words"
	"testing"
)

func statusToString(status LetterStatus) string {
	switch status {
	case None:
		return "none"
	case Correct:
		return "correct"
	case Present:
		return "present"
	case Absent:
		return "absent"
	default:
		return "unknown"
	}
}
func TestNewWordleState(t *testing.T) {
	word := "HELLO"
	ws := NewWordleState(word)

	wordString := string(ws.Word[:])
	t.Logf("Created wordleState: %+v", ws)
	t.Logf("    word: %s", wordString)
	t.Logf("    guesses: %#v", ws.guesses)
	t.Logf("    currGuess: %d", ws.currGuess)

	if wordString != word {
		t.Errorf("word = %s; want %s", wordString, word)
	}
}

func TestNewGuess(t *testing.T) {
	wordToGuess := "YIELD"
	guess := newGuess(wordToGuess)

	t.Logf("New guess: %s", guess.string())

	// Check that the letter and status are correct
	for i, l := range guess {
		t.Logf("    letter %d: %c, %s", i, l.char, statusToString(l.status))

		if l.char != wordToGuess[i] || l.status != None {
			t.Errorf(
				"letter[%d] = %c, %s; want %c, none",
				i,
				l.char,
				statusToString(l.status),
				wordToGuess[i],
			)
		}
	}
}

func TestUpdateLettersWithWord(t *testing.T) {
	guessWord := "YIELD"
	guess := newGuess(guessWord)

	var word [WordSize]byte
	copy(word[:], "HELLO")
	guess.updateLettersWithWord(word)

	statuses := []LetterStatus{
		Absent,  // "Y" is not in "HELLO"
		Absent,  // "I" is not in "HELLO"
		Present, // "E" is in "HELLO" but not in the correct position
		Correct, // "L" is in "HELLO" and in the correct position
		Absent,  // "D" is not in "HELLO"
	}

	// Check that statuses are correct
	for i, l := range guess {
		if l.status != statuses[i] {
			t.Errorf(
				"letter[%d] = %c, %s; want %c, %s",
				i,
				l.char,
				statusToString(l.status),
				guessWord[i],
				statusToString(statuses[i]),
			)
		}
	}
}

func TestAppendGuess(t *testing.T) {
	ws := NewWordleState("HELLO")

	// Test that currGuess is incremented when we add one guess
	if err := ws.appendGuess(newGuess("YIELD")); err == nil {
		if ws.currGuess != 1 {
			t.Errorf("currGuess = %d; want 1", ws.currGuess)
		}
	}

	// Add five more guesses (six guesses total)
	for i := 0; i < 5; i++ {
		guess := newGuess(words.GetWord())
		if err := ws.appendGuess(guess); err != nil {
			t.Errorf("newGuess() returned an error: %v", err)
		}
	}
}

func TestAppendGuessError(t *testing.T) {
	ws := NewWordleState("HELLO")

	var g guess
	// Add six guesses
	for i := 0; i < 6; i++ {
		g = newGuess(words.GetWord())
		ws.appendGuess(g)
	}
	// Test that we get an error when we try to add a seventh guess
	g = newGuess(words.GetWord())
	if err := ws.appendGuess(g); err == nil {
		t.Errorf("appendGuess() should result in an error")
	}
}

func TestIsWordGuessed(t *testing.T) {
	ws := NewWordleState("HELLO")
	g := newGuess("HELLO")

	g.updateLettersWithWord(ws.Word)
	ws.appendGuess(g)

	if !ws.isWordGuessed() {
		t.Errorf("isWordGuessed() should return true")
	}
}

func TestShouldEndGameCorrect(t *testing.T) {
	ws := NewWordleState("HELLO")
	g := newGuess("HELLO")

	g.updateLettersWithWord(ws.Word)
	ws.appendGuess(g)

	if !ws.shouldEndGame() {
		t.Errorf("shouldEndGame() should return true")
	}
}

func TestShouldEndGameMaxGuesses(t *testing.T) {
	ws := NewWordleState("HELLO")

	// Test that we end the game when we have six guesses
	for i := 0; i < 6; i++ {
		g := newGuess(words.GetWord())
		g.updateLettersWithWord(ws.Word)
		ws.appendGuess(g)
	}
	if !ws.shouldEndGame() {
		t.Errorf("shouldEndGame() should return true")
	}
}
