// Copyright (C) 2015 Adriano Soares <adrianosoaresjn@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type CardData struct {
	Id        int
	Ot        int
	Alias     int
	Setcode   int64
	Type      int
	Atk       int
	Def       int
	Level     int
	Race      int
	Attribute int
	Category  int64
}

type CardText struct {
	Id   int
	Name string

	Desc, Str1, Str2, Str3, Str4, Str5, Str6, Str7, Str8  string
	Str9, Str10, Str11, Str12, Str13, Str14, Str15, Str16 string
}

var dbName = flag.String("db", "cards.cdb", "database name")

// print banner
// for card in cards:
//   if card in whitelist:
//      write card to file
// print footer

func main() {
	flag.Parse()
	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)

	db := sqlx.MustOpen("sqlite3", *dbName)

	cardData := make(map[int]CardData)
	var tmp []CardData
	err := db.Select(&tmp, "select * from datas")
	catch(err)
	for _, card := range tmp {
		cardData[card.Id] = card
	}

	cardText := make(map[int]CardText)
	var tmp2 []CardText
	err = db.Select(&tmp2, "select * from texts")
	catch(err)
	for _, card := range tmp2 {
		cardText[card.Id] = card
	}

	file, err := os.Create(outputFilename)
	catch(err)
	defer file.Close()

	cardList := loadCardList(inputFilename)
	fmt.Println(len(cardList), "cards")

	_, err = file.WriteString("BEGIN TRANSACTION;\r\n")
	catch(err)

	for _, id := range cardList {
		data, found := cardData[id]
		if !found {
			file.WriteString(fmt.Sprintf("-- %d\r\n", id))
			continue
		}
		_, err := file.WriteString(formatCardData(data))
		text := cardText[id]
		_, err = file.WriteString(formatCardText(text))
		catch(err)
	}
	file.WriteString("COMMIT;\r\n")
}

func formatCardData(c CardData) string {
	return fmt.Sprintf("INSERT OR REPLACE INTO \"datas\" VALUES "+
		"(%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d);\r\n",
		c.Id, c.Ot, c.Alias, c.Setcode, c.Type, c.Atk, c.Def,
		c.Level, c.Race, c.Attribute, c.Category)
}

func formatCardText(c CardText) string {
	return fmt.Sprintf("INSERT OR REPLACE INTO \"texts\" VALUES "+
		"(%d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', "+
		"'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\r\n",
		c.Id, q(c.Name), q(c.Desc), q(c.Str1), q(c.Str2), q(c.Str3), q(c.Str4),
		q(c.Str5), q(c.Str6), q(c.Str7), q(c.Str8), q(c.Str9), q(c.Str10),
		q(c.Str11), q(c.Str12), q(c.Str13), q(c.Str14), q(c.Str15), q(c.Str16))

}

func q(s string) string {
	return strings.Replace(s, `'`, `''`, -1)
}

func loadCardList(filename string) []int {
	file, err := os.Open(filename)
	catch(err)
	defer file.Close()

	var cardList []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err == nil {
			cardList = append(cardList, id)
		}
	}
	catch(scanner.Err())
	return cardList
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
