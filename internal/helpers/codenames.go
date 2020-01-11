package helpers

import (
	"fmt"
	"math/rand"
	"strings"
)

var (
	left = [...]string{
		"admiring",
		"adoring",
		"affectionate",
		"agitated",
		"amazing",
		"angry",
		"awesome",
		"beautiful",
		"blissful",
		"bold",
		"boring",
		"brave",
		"busy",
		"charming",
		"clever",
		"cool",
		"compassionate",
		"competent",
		"condescending",
		"confident",
		"cranky",
		"crazy",
		"dazzling",
		"determined",
		"distracted",
		"dreamy",
		"eager",
		"ecstatic",
		"elastic",
		"elated",
		"elegant",
		"eloquent",
		"epic",
		"exciting",
		"fervent",
		"festive",
		"flamboyant",
		"focused",
		"friendly",
		"frosty",
		"funny",
		"gallant",
		"gifted",
		"goofy",
		"gracious",
		"great",
		"happy",
		"hardcore",
		"heuristic",
		"hopeful",
		"hungry",
		"infallible",
		"inspiring",
		"interesting",
		"intelligent",
		"jolly",
		"jovial",
		"keen",
		"kind",
		"laughing",
		"loving",
		"lucid",
		"magical",
		"mystifying",
		"modest",
		"musing",
		"naughty",
		"nervous",
		"nice",
		"nifty",
		"nostalgic",
		"objective",
		"optimistic",
		"peaceful",
		"pedantic",
		"pensive",
		"practical",
		"priceless",
		"quirky",
		"quizzical",
		"recursing",
		"relaxed",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"stoic",
		"strange",
		"stupefied",
		"suspicious",
		"sweet",
		"tender",
		"thirsty",
		"trusting",
		"unruffled",
		"upbeat",
		"vibrant",
		"vigilant",
		"vigorous",
		"wizardly",
		"wonderful",
		"xenodochial",
		"youthful",
		"zealous",
		"zen",
	}

	right = [...]string{
		"4-LOM",
		"Aayla Secura",
		"Admiral Ackbar",
		"Admiral Thrawn",
		"Ahsoka Tano",
		"Anakin Solo",
		"Asajj Ventress",
		"Aurra Sing",
		"Senator Bail Organa",
		"Barriss Offee",
		"Bastila Shan",
		"Ben Skywalker",
		"Bib Fortuna",
		"Biggs Darklighter",
		"Boba Fett",
		"Bossk",
		"Brakiss",
		"C-3PO",
		"Cad Bane",
		"Cade Skywalker",
		"Callista Ming",
		"Captain Rex",
		"Carnor Jax",
		"Chewbacca",
		"Clone Commander Cody",
		"Count Dooku",
		"Darth Bane",
		"Darth Krayt",
		"Darth Maul",
		"Darth Nihilus",
		"Darth Vader",
		"Dash Rendar",
		"Dengar",
		"Durge",
		"Emperor Palpatine",
		"Exar Kun",
		"Galen Marek",
		"General Crix Madine",
		"General Dodonna",
		"General Grievous",
		"General Veers",
		"Gilad Pellaeon",
		"Grand Moff Tarkin",
		"Greedo",
		"Han Solo",
		"IG 88",
		"Jabba The Hutt",
		"Jacen Solo",
		"Jaina Solo",
		"Jango Fett",
		"Jarael",
		"Jerec",
		"Joruus C'Baoth",
		"Ki-Adi-Mundi",
		"Kir Kanos",
		"Kit Fisto",
		"Kyle Katarn",
		"Kyp Durron",
		"Lando Calrissian",
		"Luke Skywalker",
		"Luminara Unduli",
		"Lumiya",
		"Mace Windu",
		"Mara Jade",
		"Mission Vao",
		"Natasi Daala",
		"Nom Anor",
		"Obi-Wan Kenobi",
		"Padmé Amidala",
		"Plo Koon",
		"Pre Vizsla",
		"Prince Xizor",
		"Princess Leia",
		"PROXY",
		"Qui-Gon Jinn",
		"Quinlan Vos",
		"R2-D2",
		"Rahm Kota",
		"Revan",
		"Satele Shan",
		"Savage Opress",
		"Sebulba",
		"Shaak Ti",
		"Shmi Skywalker",
		"Talon Karrde",
		"Ulic Qel-Droma",
		"Visas Marr",
		"Watto",
		"Wedge Antilles",
		"Yoda",
		"Zam Wesell",
		"Zayne Carrick",
		"Zuckuss",
	}
)

// GetRandomName generates a random name from the list of adjectives and names in this package
func GetCodename() string {
	return fmt.Sprintf("%s %s", strings.Title(left[rand.Intn(len(left))]), right[rand.Intn(len(right))])
}
