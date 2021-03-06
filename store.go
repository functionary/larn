package main

var dndcount, dnditm = 0, 0

/* number of items in the dnd inventory table	 */
const MAXITM = 83

/* this is the data for the stuff in the dnd store	 */
var itm = [90]_itm{
	/*
	 * cost 		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{2, OLEATHER, 0, 3},
	{10, OSTUDLEATHER, 0, 2},
	{40, ORING, 0, 2},
	{85, OCHAIN, 0, 2},
	{220, OSPLINT, 0, 1},
	{400, OPLATE, 0, 1},
	{900, OPLATEARMOR, 0, 1},
	{2600, OSSPLATE, 0, 1},
	{150, OSHIELD, 0, 1},

	/*
	 * cost	 	iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{2, ODAGGER, 0, 3},
	{20, OSPEAR, 0, 3},
	{80, OFLAIL, 0, 2},
	{150, OBATTLEAXE, 0, 2},
	{450, OLONGSWORD, 0, 2},
	{1000, O2SWORD, 0, 2},
	{5000, OSWORD, 0, 1},
	{16500, OLANCE, 0, 1},
	{6000, OSWORDofSLASHING, 0, 0},
	{10000, OHAMMER, 0, 0},

	/*
	 * cost		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{150, OPROTRING, 1, 1},
	{85, OSTRRING, 1, 1},
	{120, ODEXRING, 1, 1},
	{120, OCLEVERRING, 1, 1},
	{180, OENERGYRING, 0, 1},
	{125, ODAMRING, 0, 1},
	{220, OREGENRING, 0, 1},
	{1000, ORINGOFEXTRA, 0, 1},

	{280, OBELT, 0, 1},

	{400, OAMULET, 0, 1},

	{6500, OORBOFDRAGON, 0, 0},
	{5500, OSPIRITSCARAB, 0, 0},
	{5000, OCUBEofUNDEAD, 0, 0},
	{6000, ONOTHEFT, 0, 0},

	{590, OCHEST, 6, 1},
	{200, OBOOK, 8, 1},
	{10, OCOOKIE, 0, 3},

	/*
	 * cost		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{20, OPOTION, 0, 6},
	{90, OPOTION, 1, 5},
	{520, OPOTION, 2, 1},
	{100, OPOTION, 3, 2},
	{50, OPOTION, 4, 2},
	{150, OPOTION, 5, 2},
	{70, OPOTION, 6, 1},
	{30, OPOTION, 7, 7},
	{200, OPOTION, 8, 1},
	{50, OPOTION, 9, 1},
	{80, OPOTION, 10, 1},

	/*
	 * cost		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{30, OPOTION, 11, 3},
	{20, OPOTION, 12, 5},
	{40, OPOTION, 13, 3},
	{35, OPOTION, 14, 2},
	{520, OPOTION, 15, 1},
	{90, OPOTION, 16, 2},
	{200, OPOTION, 17, 2},
	{220, OPOTION, 18, 4},
	{80, OPOTION, 19, 6},
	{370, OPOTION, 20, 3},
	{50, OPOTION, 22, 1},
	{150, OPOTION, 23, 3},

	/*
	 * cost		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{100, OSCROLL, 0, 2},
	{125, OSCROLL, 1, 2},
	{60, OSCROLL, 2, 4},
	{10, OSCROLL, 3, 4},
	{100, OSCROLL, 4, 3},
	{200, OSCROLL, 5, 2},
	{110, OSCROLL, 6, 1},
	{500, OSCROLL, 7, 2},
	{200, OSCROLL, 8, 2},
	{250, OSCROLL, 9, 4},
	{20, OSCROLL, 10, 5},
	{30, OSCROLL, 11, 3},

	/*
	 * cost 		iven name		iven arg   how gp
	 * iven[]		ivenarg[]  many
	 */

	{340, OSCROLL, 12, 1},
	{340, OSCROLL, 13, 1},
	{300, OSCROLL, 14, 2},
	{400, OSCROLL, 15, 2},
	{500, OSCROLL, 16, 2},
	{1000, OSCROLL, 17, 1},
	{500, OSCROLL, 18, 1},
	{340, OSCROLL, 19, 2},
	{220, OSCROLL, 20, 3},
	{3900, OSCROLL, 21, 0},
	{610, OSCROLL, 22, 1},
	{3000, OSCROLL, 23, 0},
}

/*
	function for the dnd store
*/
func dnd_2hed() {
	lprcat("Welcome to the Larn Thrift Shoppe.  We stock many items explorers find useful\n")
	lprcat(" in their adventures.  Feel free to browse to your hearts content.\n")
	lprcat("Also be advised, if you break 'em, you pay for 'em.")
}

func dnd_hed() {
	for i := dnditm; i < 26+dnditm; i++ {
		dnditem(i)
	}
	cursor(50, 18)
	lprcat("You have ")
}

func handsfull() {
	lprcat("\nYou can't carry anything more!")
	lflush()
	nap(2200)
}

func outofstock() {
	lprcat("\nSorry, but we are out of that item.")
	lflush()
	nap(2200)
}

func nogold() {
	lprcat("\nYou don't have enough gold to pay for that!")
	lflush()
	nap(2200)
}

func dndstore() {
	dnditm = 0
	nosignal = true /* disable signals */
	clear()
	dnd_2hed()
	if outstanding_taxes > 0 {
		lprcat("\n\nThe Larn Revenue Service has ordered us to not do business with tax evaders.\n")
		beep()
		lprintf("They have also told us that you owe %d gp in back taxes, and as we must\n", outstanding_taxes)
		lprcat("comply with the law, we cannot serve you at this time.  Soo Sorry.\n")
		cursors()
		lprcat("\nPress ")
		standout("escape")
		lprcat(" to leave: ")
		lflush()
		i := 0
		for i != '\033' {
			i = ttgetch()
		}
		drawscreen()
		nosignal = false /* enable signals */
		return
	}
	dnd_hed()
	for {
		cursor(59, 18)
		lprintf("%d gold pieces", c[GOLD])
		cltoeoln()
		cl_dn(1, 20) /* erase to eod */
		lprcat("\nEnter your transaction [")
		standout("space")
		lprcat(" for more, ")
		standout("escape")
		lprcat(" to leave]? ")
		i := 0
		for (i < 'a' || i > 'z') && i != ' ' && i != '\033' && i != 12 {
			i = ttgetch()
		}
		if i == 12 {
			clear()
			dnd_2hed()
			dnd_hed()
		} else if i == '\033' {
			drawscreen()
			nosignal = false /* enable signals */
			return
		} else if i == ' ' {
			cl_dn(1, 4)
			dnditm += 26
			if dnditm >= MAXITM {
				dnditm = 0
			}
			dnd_hed()
		} else { /* buy something */
			lprc(byte(i)) /* echo the byte */
			i += dnditm - 'a'
			if i >= MAXITM {
				outofstock()
			} else if itm[i].qty <= 0 {
				outofstock()
			} else if pocketfull() {
				handsfull()
			} else if c[GOLD] < itm[i].price*10 {
				nogold()
			} else {
				if itm[i].obj == OPOTION {
					potionname[itm[i].arg] = potionhide[itm[i].arg]
				} else if itm[i].obj == OSCROLL {
					scrollname[itm[i].arg] = scrollhide[itm[i].arg]
				}
				c[GOLD] -= itm[i].price * 10
				itm[i].qty--
				take(itm[i].obj, itm[i].arg)
				if itm[i].qty == 0 {
					dnditem(i)
				}
				nap(1001)
			}
		}

	}
}

/*
	dnditem(index)

	to print the item list;  used in dndstore() enter with the index into itm
*/
func dnditem(i int) {
	if i >= MAXITM {
		return
	}
	j := (i&1)*40 + 1
	k := ((i % 26) >> 1) + 5
	cursor(j, k)
	if itm[i].qty == 0 {
		lprintf("%39s", "")
		return
	}
	lprintf("%c) ", (i%26)+'a')
	if itm[i].obj == OPOTION {
		lprintf("potion of%s", potionhide[itm[i].arg])
	} else if itm[i].obj == OSCROLL {
		lprintf("scroll of%s", scrollhide[itm[i].arg])
	} else {
		lprintf("%s", objectname[itm[i].obj])
	}
	cursor(j+31, k)
	lprintf("%6d", (itm[i].price * 10))
}

/*
	for the college of larn
*/
var course [26]byte /* the list of courses taken	 */
var coursetime = [...]int{10, 15, 10, 20, 10, 10, 10, 5}

/*
	function to display the header info for the school
*/
func sch_hed() {
	clear()
	lprcat("The College of Larn offers the exciting opportunity of higher education to\n")
	lprcat("all inhabitants of the caves.  Here is a list of the class schedule:\n\n\n")
	lprcat("\t\t    Course Name \t       Time Needed\n\n")

	if course[0] == 0 {
		lprcat("\t\ta)  Fighters Training I         10 mobuls") /* line 7 of crt */
	}
	lprc('\n')
	if course[1] == 0 {
		lprcat("\t\tb)  Fighters Training II        15 mobuls")
	}
	lprc('\n')
	if course[2] == 0 {
		lprcat("\t\tc)  Introduction to Wizardry    10 mobuls")
	}
	lprc('\n')
	if course[3] == 0 {
		lprcat("\t\td)  Applied Wizardry            20 mobuls")
	}
	lprc('\n')
	if course[4] == 0 {
		lprcat("\t\te)  Behavioral Psychology       10 mobuls")
	}
	lprc('\n')
	if course[5] == 0 {
		lprcat("\t\tf)  Faith for Today             10 mobuls")
	}
	lprc('\n')
	if course[6] == 0 {
		lprcat("\t\tg)  Contemporary Dance          10 mobuls")
	}
	lprc('\n')
	if course[7] == 0 {
		lprcat("\t\th)  History of Larn              5 mobuls")
	}

	lprcat("\n\n\t\tAll courses cost 250 gold pieces.")
	cursor(30, 18)
	lprcat("You are presently carrying ")
}

func oschool() {
	var time_used int
	nosignal = true /* disable signals */
	sch_hed()
	for {
		cursor(57, 18)
		lprintf("%d gold pieces.   ", c[GOLD])
		cursors()
		lprcat("\nWhat is your choice [")
		standout("escape")
		lprcat(" to leave] ? ")
		yrepcount = 0
		i := 0
		for (i < 'a' || i > 'h') && i != '\033' && i != 12 {
			i = ttgetch()
		}
		if i == 12 {
			sch_hed()
			continue
		} else if i == '\033' {
			nosignal = false
			drawscreen() /* enable signals */
			return
		}
		lprc(byte(i))
		if c[GOLD] < 250 {
			nogold()
		} else if course[i-'a'] != 0 {
			lprcat("\nSorry, but that class is filled.")
			nap(1000)
		} else if i <= 'h' {
			c[GOLD] -= 250
			time_used = 0
			switch i {
			case 'a':
				c[STRENGTH] += 2
				c[CONSTITUTION]++
				lprcat("\nYou feel stronger!")
				cl_line(16, 7)

			case 'b':
				if course[0] == 0 {
					lprcat("\nSorry, but this class has a prerequisite of Fighters Training I")
					c[GOLD] += 250
					time_used = -10000
					break
				}
				lprcat("\nYou feel much stronger!")
				cl_line(16, 8)
				c[STRENGTH] += 2
				c[CONSTITUTION] += 2

			case 'c':
				c[INTELLIGENCE] += 2
				lprcat("\nThe task before you now seems more attainable!")
				cl_line(16, 9)

			case 'd':
				if course[2] == 0 {
					lprcat("\nSorry, but this class has a prerequisite of Introduction to Wizardry")
					c[GOLD] += 250
					time_used = -10000
					break
				}
				lprcat("\nThe task before you now seems very attainable!")
				cl_line(16, 10)
				c[INTELLIGENCE] += 2

			case 'e':
				c[CHARISMA] += 3
				lprcat("\nYou now feel like a born leader!")
				cl_line(16, 11)

			case 'f':
				c[WISDOM] += 2
				lprcat("\nYou now feel more confident that you can find the potion in time!")
				cl_line(16, 12)

			case 'g':
				c[DEXTERITY] += 3
				lprcat("\nYou feel like dancing!")
				cl_line(16, 13)

			case 'h':
				c[INTELLIGENCE]++
				lprcat("\nYour instructor told you that the Eye of Larn is rumored to be guarded\n")
				lprcat("by a platinum dragon who possesses psionic abilities. ")
				cl_line(16, 14)
			}
			time_used += coursetime[i-'a'] * 100
			if time_used > 0 {
				gltime += time_used
				course[i-'a']++ /* remember that he has taken that course	 */
				c[HP] = c[HPMAX]
				c[SPELLS] = c[SPELLMAX] /* he regenerated */

				if c[BLINDCOUNT] != 0 {
					c[BLINDCOUNT] = 1 /* cure blindness too!  */
				}
				if c[CONFUSE] != 0 {
					c[CONFUSE] = 1 /* end confusion	 */
				}
				adjusttime(time_used) /* adjust parameters for time change */
			}
			nap(1000)
		}
	}
}

/*
 *	for the first national bank of Larn
 */
var lasttime = 0 /* last time he was in bank */

func obank() {
	banktitle("    Welcome to the First National Bank of Larn.")
}

func obank2() {
	banktitle("Welcome to the 5th level branch office of the First National Bank of Larn.")
}

func banktitle(str string) {
	nosignal = true /* disable signals */
	clear()
	lprcat(str)
	if outstanding_taxes > 0 {
		lprcat("\n\nThe Larn Revenue Service has ordered that your account be frozen until all\n")
		beep()
		lprintf("levied taxes have been paid.  They have also told us that you owe %d gp in\n", outstanding_taxes)
		lprcat("taxes, and we must comply with them. We cannot serve you at this time.  Sorry.\n")
		lprcat("We suggest you go to the LRS office and pay your taxes.\n")
		cursors()
		lprcat("\nPress ")
		standout("escape")
		lprcat(" to leave: ")
		lflush()
		i := 0
		for i != '\033' {
			i = ttgetch()
		}
		drawscreen()
		nosignal = false /* enable signals */
		return
	}
	lprcat("\n\n\tGemstone\t      Appraisal\t\tGemstone\t      Appraisal")
	obanksub()
	nosignal = false /* enable signals */
	drawscreen()
}

/*
 *	function to put interest on your bank account
 */
func ointerest() {
	if c[BANKACCOUNT] < 0 {
		c[BANKACCOUNT] = 0
	} else if c[BANKACCOUNT] > 0 && c[BANKACCOUNT] < 500000 {
		i := (gltime - lasttime) / 100 /* # mobuls elapsed */
		for {
			i--
			if i < 0 || c[BANKACCOUNT] >= 500000 {
				break
			}
			c[BANKACCOUNT] += c[BANKACCOUNT] / 250
		}
		if c[BANKACCOUNT] > 500000 {
			c[BANKACCOUNT] = 500000 /* interest limit */
		}
	}
	lasttime = (gltime / 100) * 100
}

var gemorder [26]int /* the reference to screen location for each */
var gemvalue [26]int /* the appraisal of the gems */

func obanksub() {
	ointerest() /* credit any needed interest */

	for k, i := 0, 0; i < 26; i++ {
		switch iven[i] {
		case OLARNEYE, ODIAMOND, OEMERALD, ORUBY, OSAPPHIRE:

			if iven[i] == OLARNEYE {
				gemvalue[i] = 250000 - ((gltime*7)/100)*100
				if gemvalue[i] < 50000 {
					gemvalue[i] = 50000
				}
			} else {
				gemvalue[i] = (255 & ivenarg[i]) * 100
			}
			gemorder[i] = k
			cursor((k%2)*40+1, (k>>1)+4)
			lprintf("%c) %s", i+'a', objectname[iven[i]])
			cursor((k%2)*40+33, (k>>1)+4)
			lprintf("%5d", gemvalue[i])
			k++
		default:
			// Don't allow player to sell non-existent gems
			gemvalue[i] = 0
			gemorder[i] = 0
		}
	}
	cursor(31, 17)
	lprintf("You have %8d gold pieces in the bank.", c[BANKACCOUNT])
	cursor(40, 18)
	lprintf("You have %8d gold pieces", c[GOLD])
	if c[BANKACCOUNT]+c[GOLD] >= 500000 {
		lprcat("\nNote:  Larndom law states that only deposits under 500,000gp  can earn interest.")
	}
	for {
		cl_dn(1, 20)
		lprcat("\nYour wish? [(")
		standout("d")
		lprcat(") deposit, (")
		standout("w")
		lprcat(") withdraw, (")
		standout("s")
		lprcat(") sell a stone, or ")
		standout("escape")
		lprcat("]  ")
		yrepcount = 0
		i := 0
		for i != 'd' && i != 'w' && i != 's' && i != '\033' {
			i = ttgetch()
		}
		switch i {
		case 'd':
			lprcat("deposit\nHow much? ")
			amt := readnum(c[GOLD])
			if amt < 0 {
				lprcat("\nSorry, but we can't take negative gold!")
				nap(2000)
				amt = 0
			} else if amt > c[GOLD] {
				lprcat("  You don't have that much.")
				nap(2000)
			} else {
				c[GOLD] -= amt
				c[BANKACCOUNT] += amt
			}

		case 'w':
			lprcat("withdraw\nHow much? ")
			amt := readnum(c[BANKACCOUNT])
			if amt < 0 {
				lprcat("\nSorry, but we don't have any negative gold!")
				nap(2000)
				amt = 0
			} else if amt > c[BANKACCOUNT] {
				lprcat("\nYou don't have that much in the bank!")
				nap(2000)
			} else {
				c[GOLD] += amt
				c[BANKACCOUNT] -= amt
			}

		case 's':
			lprcat("\nWhich stone would you like to sell? ")
			i := 0
			for (i < 'a' || i > 'z') && i != '*' {
				i = ttgetch()
			}
			if i == '*' {
				for i = 0; i < 26; i++ {
					if gemvalue[i] != 0 {
						c[GOLD] += gemvalue[i]
						iven[i] = 0
						gemvalue[i] = 0
						k = gemorder[i]
						cursor((k%2)*40+1, (k>>1)+4)
						lprintf("%39s", "")
					}
				}
			} else {
				i -= 'a'
				if gemvalue[i] == 0 {
					lprintf("\nItem %c is not a gemstone!", i+'a')
					nap(2000)
					break
				}
				c[GOLD] += gemvalue[i]
				iven[i] = 0
				gemvalue[i] = 0
				k = gemorder[i]
				cursor((k%2)*40+1, (k>>1)+4)
				lprintf("%39s", "")
			}

		case '\033':
			return
		}
		cursor(40, 17)
		lprintf("%8d", c[BANKACCOUNT])
		cursor(49, 18)
		lprintf("%8d", c[GOLD])
	}
}

/* XXX: apparently unused */
/*
	subroutine to appraise any stone for the bank
*/
func appraise(gemstone int) {
	/*
		int    j, amt;
		for (j = 0; j < 26; j++) {
			if (iven[j] == gemstone) {
				lprintf("\nI see you have %s", objectname[gemstone]);
				if (gemstone == OLARNEYE)
					lprcat("  I must commend you.  I didn't think\nyou could get it.");
				lprcat("  Shall I appraise it for you? ");
				yrepcount = 0;
				if (getyn() == 'y') {
					lprcat("yes.\n  Just one moment please \n");
					nap(1000);
					if (gemstone == OLARNEYE) {
						amt = 250000 - ((gltime * 7) / 100) * 100;
						if (amt < 50000)
							amt = 50000;
					} else
						amt = (255 & ivenarg[j]) * 100;
					lprintf("\nI can see this is an excellent stone, It is worth %d", (long) amt);
					lprcat("\nWould you like to sell it to us? ");
					yrepcount = 0;
					if (getyn() == 'y') {
						lprcat("yes\n");
						c[GOLD] += amt;
						iven[j] = 0;
					} else
						lprcat("no thank you.\n");
					if (gemstone == OLARNEYE)
						lprcat("It is, of course, your privilege to keep the stone\n");
				} else
					lprcat("no\nO. K.\n");
			}
		}
	*/
}

/*
	function for the trading post
*/
func otradhead() {
	clear()
	lprcat("Welcome to the Larn Trading Post.  We buy items that explorers no longer find\n")
	lprcat("useful.  Since the condition of the items you bring in is not certain,\n")
	lprcat("and we incur great expense in reconditioning the items, we usually pay\n")
	lprcat("only 20% of their value were they to be new.  If the items are badly\n")
	lprcat("damaged, we will pay only 10% of their new value.\n\n")
}

func otradepost() {
	dnditm, dndcount = 0, 0
	nosignal = true /* disable signals */
	resetscroll()
	otradhead()
	for {
		lprcat("\nWhat item do you want to sell to us [")
		standout("*")
		lprcat(" for list, or ")
		standout("escape")
		lprcat("] ? ")
		i := 0
		for i > 'z' || (i < 'a' && i != '*' && i != '\033' && i != '.') {
			i = ttgetch()
		}
		if i == '\033' {
			setscroll()
			recalc()
			drawscreen()
			nosignal = false /* enable signals */
			return
		}
		isub := i - 'a'
		j := 0
		oor := false
		if isub < 0 || isub >= len(iven) {
			oor = true
		}
		if !oor && iven[isub] == OSCROLL {
			if scrollname[ivenarg[isub]] == "" {
				j = 1
				cnsitm()
			} /* can't sell unidentified item */
		}
		if !oor && iven[isub] == OPOTION {
			if potionname[ivenarg[isub]] == "" {
				j = 1
				cnsitm()
			} /* can't sell unidentified item */
		}
		if j == 0 {
			if i == '*' {
				clear()
				qshowstr()
				otradhead()
			} else if oor || iven[isub] == 0 {
				lprintf("\nYou don't have item %c!", isub+'a')
			} else {
				for j = 0; j < MAXITM; j++ {
					if itm[j].obj == iven[isub] || iven[isub] == ODIAMOND || iven[isub] == ORUBY || iven[isub] == OEMERALD || iven[isub] == OSAPPHIRE {
						srcount = 0
						show3(isub) /* show what the item was */
						var value int
						if iven[isub] == ODIAMOND || iven[isub] == ORUBY || iven[isub] == OEMERALD || iven[isub] == OSAPPHIRE {
							value = 20 * ivenarg[isub]
						} else if itm[j].obj == OSCROLL || itm[j].obj == OPOTION {
							value = 2 * itm[j+ivenarg[isub]].price
						} else {
							izarg := ivenarg[isub]
							value = itm[j].price /* appreciate if a +n object */
							if izarg >= 0 {
								value *= 2
							}
							for {
								izarg--
								if izarg < 0 {
									break
								}
								value = 14 * (67 + value) / 10
								if value >= 500000 {
									break
								}
							}
						}
						lprintf("\nItem (%c) is worth %d gold pieces to us.  Do you want to sell it? ", i, value)
						yrepcount = 0
						if getyn() == 'y' {
							lprcat("yes\n")
							c[GOLD] += value
							if c[WEAR] == isub {
								c[WEAR] = -1
							}
							if c[WIELD] == isub {
								c[WIELD] = -1
							}
							if c[SHIELD] == isub {
								c[SHIELD] = -1
							}
							adjustcvalues(iven[isub], ivenarg[isub])
							iven[isub] = 0
						} else {
							lprcat("no thanks.\n")
						}
						j = MAXITM + 100 /* get out of the inner loop */
					}
				}
				if j <= MAXITM+2 {
					lprcat("\nSo sorry, but we are not authorized to accept that item.")
				}
			}
		}
	}
}

func cnsitm() {
	lprcat("\nSorry, we can't accept unidentified objects.")
}

/*
 *	for the Larn Revenue Service
 */
func olrs() {
	nosignal = true /* disable signals */
	first := 1
	clear()
	resetscroll()
	cursor(1, 4)
	lprcat("Welcome to the Larn Revenue Service district office.  How can we help you?")
	var i int
	for {
		if first != 0 {
			first = 0
			goto nxt
		}
		cursors()
		lprcat("\n\nYour wish? [(")
		standout("p")
		lprcat(") pay taxes, or ")
		standout("escape")
		lprcat("]  ")
		yrepcount = 0
		i = 0
		for i != 'p' && i != '\033' {
			i = ttgetch()
		}
		switch i {
		case 'p':
			lprcat("pay taxes\nHow much? ")
			amt := readnum(c[GOLD])
			if amt < 0 {
				lprcat("\nSorry, but we can't take negative gold\n")
				amt = 0
			} else if amt > c[GOLD] {
				lprcat("  You don't have that much.\n")
			} else {
				c[GOLD] -= paytaxes(amt)
			}

		case '\033':
			nosignal = false /* enable signals */
			setscroll()
			drawscreen()
			return
		}

	nxt:
		cursor(1, 6)
		if outstanding_taxes > 0 {
			lprintf("You presently owe %d gp in taxes.  ", outstanding_taxes)
		} else {
			lprcat("You do not owe us any taxes.           ")
		}
		cursor(1, 8)
		if c[GOLD] > 0 {
			lprintf("You have %6d gp.    ", c[GOLD])
		} else {
			lprcat("You have no gold pieces.  ")
		}
	}
}
