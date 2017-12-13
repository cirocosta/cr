package lib

import (
	"math/rand"
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
		"blissful",
		"boring",
		"brave",
		"clever",
		"cocky",
		"compassionate",
		"competent",
		"condescending",
		"confident",
		"cranky",
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
		"fervent",
		"festive",
		"flamboyant",
		"focused",
		"friendly",
		"frosty",
		"gallant",
		"gifted",
		"goofy",
		"gracious",
		"happy",
		"hardcore",
		"heuristic",
		"hopeful",
		"hungry",
		"infallible",
		"inspiring",
		"jolly",
		"jovial",
		"keen",
		"kind",
		"laughing",
		"loving",
		"lucid",
		"mystifying",
		"modest",
		"musing",
		"naughty",
		"nervous",
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
		"relaxed",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"stoic",
		"stupefied",
		"suspicious",
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
		"albattani",

		"allen",

		"almeida",

		"agnesi",

		"archimedes",

		"ardinghelli",

		"aryabhata",

		"austin",

		"babbage",

		"banach",

		"bardeen",

		"bartik",

		"bassi",

		"beaver",

		"bell",

		"benz",

		"bhabha",

		"bhaskara",

		"blackwell",

		"bohr",

		"booth",

		"borg",

		"bose",

		"boyd",

		"brahmagupta",

		"brattain",

		"brown",

		"carson",

		"chandrasekhar",

		"chatterjee",

		"shannon",

		"clarke",

		"colden",

		"cori",

		"cray",

		"curran",

		"curie",

		"darwin",

		"davinci",

		"dijkstra",

		"dubinsky",

		"easley",

		"edison",

		"einstein",

		"elion",

		"engelbart",

		"euclid",

		"euler",

		"fermat",

		"fermi",

		"feynman",

		"franklin",

		"galileo",

		"gates",

		"goldberg",

		"goldstine",

		"goldwasser",

		"golick",

		"goodall",

		"haibt",

		"hamilton",

		"hawking",

		"heisenberg",

		"hermann",

		"heyrovsky",

		"hodgkin",

		"hoover",

		"hopper",

		"hugle",

		"hypatia",

		"jackson",

		"jang",

		"jennings",

		"jepsen",

		"johnson",

		"joliot",

		"jones",

		"kalam",

		"kare",

		"keller",

		"kepler",

		"khorana",

		"kilby",

		"kirch",

		"knuth",

		"kowalevski",

		"lalande",

		"lamarr",

		"lamport",

		"leakey",

		"leavitt",

		"lewin",

		"lichterman",

		"liskov",

		"lovelace",

		"lumiere",

		"mahavira",

		"mayer",

		"mccarthy",

		"mcclintock",

		"mclean",

		"mcnulty",

		"meitner",

		"meninsky",

		"mestorf",

		"minsky",

		"mirzakhani",

		"morse",

		"murdock",

		"neumann",

		"newton",

		"nightingale",

		"nobel",

		"noether",

		"northcutt",

		"noyce",

		"panini",

		"pare",

		"pasteur",

		"payne",

		"perlman",

		"pike",

		"poincare",

		"poitras",

		"ptolemy",

		"raman",

		"ramanujan",

		"ride",

		"montalcini",

		"ritchie",

		"roentgen",

		"rosalind",

		"saha",

		"sammet",

		"shaw",

		"shirley",

		"shockley",

		"sinoussi",

		"snyder",

		"spence",

		"stallman",

		"stonebraker",

		"swanson",

		"swartz",

		"swirles",

		"tereshkova",

		"tesla",

		"thompson",

		"torvalds",

		"turing",

		"varahamihira",

		"vaughan",

		"visvesvaraya",

		"volhard",

		"villani",

		"wescoff",

		"wiles",

		"williams",

		"wilson",

		"wing",

		"wozniak",

		"wright",

		"yalow",

		"yonath",
	}
)

// GetRandomName returns a "dockerish" random name.
// Extracted from https://github.com/moby/moby/tree/5dc791c2debd561f61f04ec5a947f261fe79b275/pkg/namesgenerator
func GetRandomName() string {
	return left[rand.Intn(len(left))] +
		"_" +
		right[rand.Intn(len(right))]
}
