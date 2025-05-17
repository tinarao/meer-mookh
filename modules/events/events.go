package events

import "meermookh/modules/player"

func PlayerKilledEnemy(p *player.Player) {
	p.AddFrags(1)
}
