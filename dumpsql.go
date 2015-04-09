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
// for card in card:
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
	db.Select(&tmp, "select * from datas")
	for _, card := range tmp {
		cardData[card.Id] = card
	}

	file, err := os.Create(outputFilename)
	catch(err)
	defer file.Close()

	cardList := loadCardList(inputFilename)
	fmt.Println(len(cardList), "loaded")

	_, err = file.WriteString("BEGIN TRANSACTION;\r\n")
	catch(err)

	for _, id := range cardList {
		c := cardData[id]
		fmt.Println("found", c)
		_, err := file.WriteString(formatCardData(c))
		catch(err)
	}
	file.WriteString("COMMIT;\r\n")
}

func formatCardData(c CardData) string {
	return fmt.Sprintf("INSERT OR REPLACE INTO \"datas\" VALUES (%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d);\r\n",
		c.Id, c.Ot, c.Alias, c.Setcode, c.Type, c.Atk, c.Def,
		c.Level, c.Race, c.Attribute, c.Category)
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
