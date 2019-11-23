package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
  "regexp"
	"net/http"
	"io/ioutil"
	"log"
	"strings"

  "github.com/gookit/color"
	"github.com/bwmarrin/discordgo"
)


var T,_ = ioutil.ReadFile("token.txt")
var Token string = strings.Replace(string(T), "\n", "", 1)

func main() {
  fmt.Println("Developed by sleep")
  dg, err := discordgo.New(Token)
	_, err = dg.User("@me")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func redeem_code(code string) {
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://discordapp.com/api/v6/entitlements/gift-codes/"+code+"/redeem", nil)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("Authorization", Token)
		resp, err := client.Do(req)
		defer resp.Body.Close()

		//bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    /*bodyString := string(bodyBytes)
		return bodyString*/
		switch status := resp.StatusCode; status {
		case 200:
			color.Success.Println("Redeemed Nitro: " + code)
		case 404:
			color.Danger.Println("Invalid Code: " + code)
		case 400:
			color.Danger.Println("Invalid Code: " + code)
		//case 400:
			//color.Comment.Println("Code Already Used: " + code)
		}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

  match1, _ := regexp.MatchString("(?:https?:)?discord.gift.(\\S+)", m.Content)
  if match1 == true {
    re := regexp.MustCompile("(?:https?:)?discord.gift.(\\S+)")
    match := re.FindStringSubmatch(m.Content)
    code := match[1]
    redeem_code(code)
}
}