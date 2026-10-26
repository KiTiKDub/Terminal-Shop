package main

import (
	"strconv"
)

func openPluginView(m model) string {

	var msg string

	switch m.choice {
	case "Disburser":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Disburser.jpg", m.list.Width()*3) //change width calculation to be straignt from teamsg eventually
	case "QwikRef":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("QwikRef.png", m.list.Width()*3)
	case "ChopChop!":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("chopchop.jpg", m.list.Width()*3)
	case "Mangl3r":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Mangl3r.jpg", m.list.Width()*3)
	case "Distork":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("Distork.png", m.list.Width()*3) + "\n" + strconv.Itoa(m.list.Width())
	case "Krush":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("KITIK_LOGO_NO_BKGD.png", m.list.Width()*3)
	case "KiTiK Kollection Vol 1":
		msg = "Price: FREE\n\n" + renderANSIHalfBlocks("kol.jpg", m.list.Width()*3)
	default:
		msg = "default" + m.choice
	}

	return msg
}
