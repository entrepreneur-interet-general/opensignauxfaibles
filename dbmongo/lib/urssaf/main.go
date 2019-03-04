package urssaf

import (
	"dbmongo/lib/engine"
)

type parserFunc func(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event)

// Parser fournit le contenu des fichiers urssaf
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "urssafParser",
		Channel: eventChannel,
	}

	go func() {
		defer close(outputChannel)
		defer close(eventChannel)

		functions := []parserFunc{
			parseCCSF,
			parseCotisation,
      parseCompte,
			parseDebit,
			parseDelai,
			parseDPAE,
			parseEffectif,
			parseProcol,
		}

		mapping, err := getCompteSiretMapping(&batch)
		if err != nil {
			event.Critical("Erreur mapping: " + err.Error())
			return
		}

		for _, f := range functions {
			outputs, events := f(batch, mapping)
			go engine.PlugEvents(events, eventChannel)
			for o := range outputs {
				outputChannel <- o
			}
		}

	}()
	return outputChannel, eventChannel
}
