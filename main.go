package main

import (
	"fmt"
	//www.github.com/urmel/go2d
	"go2d"
	"rand"

	"json"
	"io/ioutil"
	"time"
	//"os"
)

const (
	SCREEN_WIDTH = 800
	SCREEN_HEIGHT = 600
)

var ticker uint32 = go2d.GetTicks()
var humans []*Human = make([]*Human, 0)
var font *go2d.Font
var borden [5]*go2d.Image
var hatchOpen bool = false

var temp_head int = 0
var temp_upper int = 0
var temp_lower int = 0

var epidemics [4]string

var currentCountry *Country 
var currentDiseases []*Disease
var currentEnv *Env

type Disease struct {
	name string
	head, upper, lower bool
	epidemic bool
}

type Country struct {
	name string
	diseases []string
}

func NewCountry(name string) *Country {
	country := &Country{}
	country.name = name
	country.diseases = make([]string, 0)
	return country
}

func (country *Country) addDisease(name string) {
	country.diseases = append(country.diseases, name)
}

var countries []*Country
var diseases map[string]*Disease = make(map[string]*Disease)

func setDiseases() {
	currentDiseases = make([]*Disease, 0)
	var head, upper, lower bool
	if temp_head > 37 {
		head = true
	}
	if temp_upper > 37 {
		upper = true
	}
	if temp_lower > 37 {
		lower = true
	}
	if currentCountry != nil {
		for _, diseasename := range currentCountry.diseases {
			disease := diseases[diseasename]
			if disease == nil {
				continue
			}
			
			scouterTrue := 0
			if head {
				scouterTrue++
			}
			if upper {
				scouterTrue++
			}
			if lower {
				scouterTrue++
			}
			
			diseaseTrue := 0
			if disease.head {
				diseaseTrue ++
			}
			if disease.upper {
				diseaseTrue ++
			}
			if disease.lower {
				diseaseTrue ++
			}
			
			headCheck := (head == disease.head)
			upperCheck := (upper == disease.upper)
			lowerCheck := (lower == disease.lower)
			
			if scouterTrue > diseaseTrue {
				if headCheck || upperCheck || lowerCheck {
					currentDiseases = append(currentDiseases, disease)
				}
			} else {
				if headCheck && upperCheck && lowerCheck {
					currentDiseases = append(currentDiseases, disease)
				}
			}
		}
	}
}

type Human struct {
	x, y int
	frames []*go2d.Image
	frame int
}

type Env struct{
	seaker *go2d.Image
	customs * go2d.Image
	bg *go2d.Image
	entrance * go2d.Image
	hatch * go2d.Image
	hatch_open * go2d.Image
	pole * go2d.Image
}

func (human *Human) addFrame(image *go2d.Image) {
	human.frames = append(human.frames, image)
}

func NewHuman() *Human {
	human := &Human{y : 400}
	human.frames = make([]*go2d.Image, 0)
	
	human.addFrame(go2d.NewImage("chiyo1.png"))
	human.addFrame(go2d.NewImage("chiyo2.png"))
	human.addFrame(go2d.NewImage("chiyo3.png"))
	
	humans = append(humans, human)
	return human
}

func removeHuman(index int) {
	h := humans[0:index]
	l := humans[index+1:]
	humans = append(h, l...)
}

func setEpidemic(index int) {
	for i := 0; i < 4; i++ {
		for {
			rand.Seed(time.Nanoseconds())
			diseasenum := rand.Intn(len(diseases))
			curnum := 0
			found := false
			var diseasename string
			for diseasename, _ = range diseases {
				if curnum == diseasenum {
					for _, dname := range epidemics {
						if diseasename == dname {
							found = true
						}
					}
					break
				}
				curnum++
			}
			if !found {
				epidemics[index] = diseasename
				return
			}
		}
	}
}

func start() {
	loadData()
	
	//set epidemics
	for i := 0; i < len(epidemics); i++ {
		setEpidemic(i)
		println(fmt.Sprintf("Epidemic #%d: %s", i, epidemics[i]))
	}
	
	borden[0] = go2d.NewImage("bord_leeg.png")
	for i := 1; i <= 4; i++ {
		borden[i] = go2d.NewImage(fmt.Sprintf("bord%d.png", i))
	}	
	
	currentEnv = &Env{}
	currentEnv.seaker = go2d.NewImage("seaker.png")
	currentEnv.bg = go2d.NewImage("bg.png")
	currentEnv.entrance = go2d.NewImage("entrance.png")
	currentEnv.customs = go2d.NewImage("desk.png")
	currentEnv.hatch = go2d.NewImage("hatch.png")
	currentEnv.hatch_open = go2d.NewImage("hatch_open.png")
	currentEnv.pole = go2d.NewImage("pole.png")
	
	font = go2d.NewFont("arial.ttf", 14)
	font.SetStyle(true, false, false)
	NewHuman()
}

func update() {
	if (go2d.GetTicks()-ticker) >= 100 {
		for i := 0; i < len(humans); i++ {
			if humans[i].x >= 500 && humans[i].x <= 600 && hatchOpen {
				humans[i].y += 20
				if humans[i].y >= SCREEN_HEIGHT {
					hatchOpen = false
					currentDiseases = make([]*Disease, 0)
					removeHuman(i)
					NewHuman()
					break
				}
				continue
			}
			if humans[i].x >= 150 && humans[i].x <= 230 {
				temp_head = 34+rand.Intn(7)
				temp_upper = 34+rand.Intn(7)
				temp_lower = 34+rand.Intn(7)
			}
			if humans[i].x >= 320 && humans[i].x <= 400 {
				country := countries[rand.Intn(len(countries))]
				if country != nil {
					currentCountry = country
					setDiseases()
				}
			}
			if humans[i].x >= 500 && humans[i].x <= 600 {
				if len(currentDiseases) > 0 {
					for _, dis := range currentDiseases {
						for _, epi := range epidemics {
							if dis.name == epi {
								hatchOpen = true
								//goto is bad? http://kerneltrap.org/node/553/2131
								goto thebatmobile
							}
						}
					}
					thebatmobile:
				}
			}
			
			humans[i].x += 10
			
			if humans[i].frame+1 < len(humans[i].frames) {
				humans[i].frame++
			} else {
				humans[i].frame = 0
			}
			
			if humans[i].x >= SCREEN_WIDTH {
				removeHuman(i)
				NewHuman()
				break
			}
		}
		ticker = go2d.GetTicks()
	}
}

func draw() {
	currentEnv.bg.DrawRect(go2d.NewRect(0,0, 800, 600))
	font.SetStyle(false, false, true)
	currentEnv.customs.DrawRect(go2d.NewRect(300, 430, 164, 82))
	
	
	if !hatchOpen {
		currentEnv.hatch.DrawRect(go2d.NewRect(470, 480, 144, 80))
	} else {
		currentEnv.hatch_open.DrawRect(go2d.NewRect(470, 480, 144, 80))	
	}
	
	
	//go2d.DrawFillRect(go2d.NewRect(320, 455, 80, 80), 255, 255, 255, 255)
	currentEnv.pole.DrawRect(go2d.NewRect(540, 280, 30, 160))
	borden[0].Draw(480, 260)
	for i := 0; i < len(humans); i++ {
		human := humans[i]
		if len(human.frames) > 0 {
			human.frames[human.frame].DrawInRect(humans[i].x, humans[i].y, go2d.NewRect(0,0, 800, 512))
		}
	}
	currentEnv.seaker.DrawRect(go2d.NewRect(150, 309, 91, 211))
	

	currentEnv.entrance.DrawRect(go2d.NewRect(717, 311, 83, 176))
	font.SetStyle(true, false, false)
	font.DrawText("Body heat:", 100, 100)
	font.DrawText("Head:", 100, 120)
	font.DrawText("Upper body:", 100, 140)
	font.DrawText("Lower body:", 100, 160)
	font.DrawText(fmt.Sprintf("%d", temp_head), 200, 120)
	font.DrawText(fmt.Sprintf("%d", temp_upper), 200, 140)
	font.DrawText(fmt.Sprintf("%d", temp_lower), 200, 160)
	
	font.DrawText("Country:", 300, 100)
	if currentCountry != nil {
		font.DrawText(currentCountry.name, 300, 120)
	}
	
	font.DrawText("Possible diseases:", 500, 100)

	if currentDiseases != nil {
		counter := 0
		for _, disease := range currentDiseases {
			checks := ""
			if disease.head {
				checks = "H"
			}
			if disease.upper {
				if checks != "" {
					checks = checks + ":U"
				} else {
					checks = "U"
				}
			}
			if disease.lower {
				if checks != "" {
					checks = checks + ":L"
				} else {
					checks = "L"
				}
			}
			font.DrawText(disease.name+" ("+checks+")", 500, 120+(counter*20))
			counter++
			
			for i, dname := range epidemics {
				if dname == disease.name {
					borden[i+1].Draw(480, 260)
					break
				}
			}
		}
	}
	
	font.SetStyle(false, false, true)
} 

func mouseMove(x, y int16) {
	//mouse move events
}

func mouseUp(x, y int16) {
	//mouse up events
}

func mouseDown(x, y int16) {
	//mouse down events
}

func textInput(char uint8) {
	//text input events
}

func keyDown(key int) {
	//key down events
}

func loadData() {
	file, _ := ioutil.ReadFile("data.json")
	var jsontype jsonobject
    json.Unmarshal(file, &jsontype)
	
	countries = make([]*Country, 0)
	for _, country := range jsontype.Data.Countries {
		newCountry := NewCountry(country.Name)
		for _, disease := range country.Diseases {
			if _, found := diseases[disease.Name]; !found {
				diseases[disease.Name] = &Disease{disease.Name, disease.Head, disease.Upper, disease.Lower, false}
			}
			newCountry.addDisease(disease.Name)
		}
		countries = append(countries, newCountry)
	}
}

////////////JSON STUFF/////////////////
type jsonobject struct {
	Data DataType
}

type DataType struct {
	Countries []CountryType
}

type CountryType struct {
	Name string
	Diseases []DiseaseType
}

type DiseaseType struct {
	Name string
	Head bool
	Upper bool
	Lower bool
}
/////////////////////////////////////////

func main() {
	game := go2d.NewGame("Heatseeker")
	game.SetDimensions(SCREEN_WIDTH, SCREEN_HEIGHT)

	//Set to false when OpenGL should also be defaulted on Windows
	game.SetD3D(true)

	game.SetInitFun(start)
	game.SetUpdateFun(update)
	game.SetDrawFun(draw)

	game.SetMouseMoveFun(mouseMove)
	game.SetMouseDownFun(mouseDown)
	game.SetMouseUpFun(mouseUp)
	game.SetTextInputFun(textInput)
	game.SetKeyDownFun(keyDown)

	game.Run()
}

