package chainserver

import (
	"net/http"

	"github.com/icfoxy/GoTools"
)

func GetChainHandler(w http.ResponseWriter, r *http.Request) {
	GoTools.RespondByJSON(w, 200, MainChain)
}

func TestAliveHandler(w http.ResponseWriter, r *http.Request) {
	GoTools.RespondByJSON(w, 200, "I AM ALIVE")
}
