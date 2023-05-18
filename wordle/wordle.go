package wordle

import (
	"errors"
	words "golang-addon/week-1/golang-clidle/words"
	"strings"
)

const (
	MaxGuesses = 6
	WordSize   = 5
)

// LetterStatus can be none, correct, present, or absent
type LetterStatus int

const (
	// None = no status, not guessed yet
	None LetterStatus = iota
	// Absent = not in the word
	Absent
	// Present = in the word, but not in the correct position
	Present
	// Correct = in the correct position
	Correct
)

type WordleState struct {
	// Word is the Word that the user is trying to guess
	Word [WordSize]byte
	// Guesses holds the Guesses that the user has made.
	// currGuess is the index of the available slot in Guesses.
	Guesses   [MaxGuesses]Guess
	CurrGuess int
}

// NewWordleState builds a new wordleState from a string.
// Pass in the word you want the user to guess.
func NewWordleState(word string) WordleState {
	w := WordleState{}
	copy(w.Word[:], word)
	return w
}

// AppendGuess adds a guess to the wordleState. It returns an error
// if the guess is invalid.
func (w *WordleState) AppendGuess(g Guess) error {
	// Reject guesses when the max number of guesses has been reached
	if w.CurrGuess >= MaxGuesses {
		return errors.New("max guesses reached")
	}

	// Reject guesses that are not the correct length
	if len(g) != WordSize {
		return errors.New("invalid guess length")
	}

	// Reject guesses that are invalid words
	if !words.IsWord(g.string()) {
		return errors.New("invalid guess word")
	}

	w.Guesses[w.CurrGuess] = g
	w.CurrGuess++
	return nil
}

// string converts a guess into a string
func (g *Guess) string() string {
	str := ""
	for _, l := range g {
		if 'A' <= l.Char && l.Char <= 'Z' {
			str += string(l.Char)
		}
	}
	return str
}

// isWordGuessed returns true when the latest guess is the correct word
func (w *WordleState) isWordGuessed() bool {
	wordIsCorrect := true
	for _, l := range w.Guesses[w.CurrGuess-1] {
		if l.Status != Correct {
			wordIsCorrect = false
		}
	}
	return wordIsCorrect
}

// shouldEndGame checks if the game should end
func (w *WordleState) shouldEndGame() bool {
	// End the game if the word is correct or the max guesses have been reached
	if w.isWordGuessed() || w.CurrGuess >= MaxGuesses {
		return true
	} else {
		return false
	}
}

type letter struct {
	Char   byte
	Status LetterStatus
}

// newLetter builds a new letter from a byte
func newLetter(char byte) letter {
	return letter{Char: char, Status: None}
}

type Guess [WordSize]letter

// NewGuess builds a new guess from a string
func NewGuess(guessedWord string) Guess {
	guess := Guess{}
	for i, c := range guessedWord {
		guess[i] = newLetter(byte(c))
	}

	return guess
}

// updateLettersWithWord updates the status of the letters in the guess based on a word
func (g *Guess) UpdateLettersWithWord(word [WordSize]byte) {
	for i := range g {
		l := &g[i]
		if l.Char == word[i] {
			l.Status = Correct
		} else if strings.Contains(string(word[:]), string(l.Char)) {
			l.Status = Present
		} else {
			l.Status = Absent
		}
	}
}
