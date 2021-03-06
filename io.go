package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/dsymonds/gocurse/curses"
)

/*
 * Below are the functions in this file:
 *
 * setupvt100() 	Subroutine to set up terminal in correct mode for game
 * clearvt100()  	Subroutine to clean up terminal when the game is over
 * ttgetch() 		Routine to read in one character from the terminal
 * scbr()			Function to set cbreak -echo for the terminal
 * sncbr()			Function to set -cbreak echo for the terminal
 * newgame() 		Subroutine to save the initial time and seed rnd()
 *
 * FILE OUTPUT ROUTINES
 *
 * lprintf(format,args . . .)	printf to the output buffer lprint(integer)
 * end binary integer to output buffer lwrite(buf,len)
 * rite a buffer to the output buffer lprcat(str)
 * ent string to output buffer
 *
 * FILE OUTPUT MACROS (in header.h)
 *
 * lprc(character)				put the character into the output
 * buffer
 *
 * FILE INPUT ROUTINES
 *
 * long lgetc()				read one character from input buffer
 * long larn_lrint()			read one integer from input buffer
 * lrfill(address,number)		put input bytes into a buffer char
 * *lgetw()				get a whitespace ended word from
 * input char *lgetl()				get a \n or EOF ended line
 * from input
 *
 * FILE OPEN / CLOSE ROUTINES
 *
 * lcreat(filename)			create a new file for write
 * lopen(filename)				open a file for read
 * lappend(filename)			open for append to an existing file
 * lrclose()					close the input file
 * lwclose()					close output file lflush()
 * lush the output buffer
 *
 * Other Routines
 *
 * cursor(x,y)					position cursor at [x,y]
 * cursors()					position cursor at [1,24]
 * (saves memory) cl_line(x,y)         		Clear line at [1,y] and leave
 * cursor at [x,y] cl_up(x,y)    				Clear screen
 * from [x,1] to current line. cl_dn(x,y)
 * lear screen from [1,y] to end of display. standout(str)
 * rint the string in standout mode. set_score_output()
 * alled when output should be literally printed. * ttputch(ch)
 * rint one character in decoded output buffer. * flush_buf()
 * lush buffer with decoded output. * init_term()
 * erminal initialization -- setup termcap info *	char *tmcapcnv(sd,ss)
 * outine to convert VT100 \33's to termcap format beep()
 * e to emit a beep if enabled (see no-beep in .larnopts)
 *
 * Note: ** entries are available only in termcap mode.
 */

/*
#ifdef TERMIO
#include <termio.h>
#define sgttyb termio
#define stty(_a,_b) ioctl(_a,TCSETA,_b)
#define gtty(_a,_b) ioctl(_a,TCGETA,_b)
#endif
#ifdef TERMIOS
#include <termios.h>
#define sgttyb termios
#define stty(_a,_b) tcsetattr(_a,TCSADRAIN,_b)
#define gtty(_a,_b) tcgetattr(_a,_b)
#endif
*/

/*
#if defined(TERMIO) || defined(TERMIOS)
static int      rawflg = 0;
static char     saveeof, saveeol;
#define doraw(_a) \
	if(!rawflg) { \
		++rawflg; \
		saveeof = _a.c_cc[VMIN]; \
		saveeol = _a.c_cc[VTIME]; \
	} \
    	_a.c_cc[VMIN] = 1; \
	_a.c_cc[VTIME] = 1; \
	_a.c_lflag &= ~(ICANON|ECHO|ECHOE|ECHOK|ECHONL)
#define unraw(_a) \
	_a.c_cc[VMIN] = saveeof; \
	_a.c_cc[VTIME] = saveeol; \
	_a.c_lflag |= ICANON|ECHO|ECHOE|ECHOK|ECHONL

#else	// not TERMIO or TERMIOS

#define CBREAK RAW		// V7 has no CBREAK

#define doraw(_a) (_a.sg_flags |= CBREAK,_a.sg_flags &= ~ECHO)
#define unraw(_a) (_a.sg_flags &= ~CBREAK,_a.sg_flags |= ECHO)
#include <sgtty.h>
#endif	// not TERMIO or TERMIOS
*/

const debugFilename = "larn-debug.log"

var (
	debug     = flag.Bool("debug", false, "whether to log debugging information to "+debugFilename)
	debugFile *os.File
)

func debugf(format string, args ...interface{}) {
	if !*debug {
		return
	}
	if debugFile == nil {
		var err error
		debugFile, err = os.Create(debugFilename)
		if err != nil {
			log.Fatalf("os.Create(%q): %v", debugFilename, err)
		}
		fmt.Fprintf(debugFile, "-----[ Larn debug file opened %v ]-----\n", time.Now())
	}
	// Walk back until we get something that doesn't look like a closure.
	for i := 1; ; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if strings.Contains(f.Name(), ".func·") {
			continue
		}
		args = append([]interface{}{path.Base(file), line, strings.TrimPrefix(f.Name(), "main.")}, args...)
		break
	}
	fmt.Fprintf(debugFile, "%s:%d\t[%s] "+format+"\n", args...)
}

const LINBUFSIZE = 128 /* size of the lgetw() and lgetl() buffer */
var io_out *os.File    /* output file number */
var io_in *os.File     /* input file */
//static struct sgttyb ttx;/* storage for the tty modes */

/*
 *	setupvt100() Subroutine to set up terminal in correct mode for game
 *
 *	Attributes off, clear screen, set scrolling region, set tty mode
 */
func setupvt100() {
	debugf("")
	clear()
	setscroll()
	scbr() /* system("stty cbreak -echo"); */
}

/*
 *	clearvt100() 	Subroutine to clean up terminal when the game is over
 *
 *	Attributes off, clear screen, unset scrolling region, restore tty mode
 */
func clearvt100() {
	debugf("")
	resetscroll()
	clear()
	sncbr() /* system("stty -cbreak echo"); */
	if err := curses.Endwin(); err != nil {
		debugf("curses.Endwin: %v", err)
	}
	win = nil
}

var win *curses.Window

/*
 *	ttgetch() 	Routine to read in one character from the terminal
 */
func ttgetch() int {
	if b, ok := popReplay(); ok {
		return int(b)
	}

	ch := win.Getch()
	debugf("win.Getch() => %d", ch)
	return ch
}

/*
 *	scbr()		Function to set cbreak -echo for the terminal
 *
 *	like: system("stty cbreak -echo")
 */
func scbr() {
	debugf("")
	if err := curses.Cbreak(); err != nil {
		debugf("curses.Cbreak: %v", err)
	}
	if err := curses.Noecho(); err != nil {
		debugf("curses.Noecho: %v", err)
	}
}

/*
 *	sncbr()		Function to set -cbreak echo for the terminal
 *
 *	like: system("stty -cbreak echo")
 */
func sncbr() {
	debugf("")
	if err := curses.Nocbreak(); err != nil {
		debugf("curses.Nocbreak: %v", err)
	}
	if err := curses.Echo(); err != nil {
		debugf("curses.Echo: %v", err)
	}
}

/*
 *	newgame() 	Subroutine to save the initial time and seed rnd()
 */
func newgame(seed uint32) {
	for i := 0; i < 100; i++ {
		c[i] = 0
	}
	seedrand(seed)
	lcreat("") /* open buffering for output to terminal */
}

/*
 *	lprintf(format,args . . .)		printf to the output buffer
 *		char *format;
 *		??? args . . .
 *
 *	Enter with the format string in "format", as per printf() usage
 *		and any needed arguments following it
 *	Note: lprintf() only supports %s, %c and %d, with width modifier and left
 *		or right justification.
 *	No correct checking for output buffer overflow is done, but flushes
 *		are done beforehand if needed.
 *	Returns nothing of value.
 */
func lprintf(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)

	if len(lpbuf) >= cap(lpbuf) {
		lflush()
	}

	lprcat(buf)
}

/*
 *	lprint(long-integer)	send binary integer to output buffer
 *		long integer;
 *
 *		+---------+---------+---------+---------+
 *		|   high  |	    |	      |	  low	|
 *		|  order  |	    |	      |  order	|
 *		|   byte  |	    |	      |	  byte	|
 *		+---------+---------+---------+---------+
 *	        31  ---  24 23 --- 16 15 ---  8 7  ---   0
 *
 *	The save order is low order first, to high order (4 bytes total)
 *	and is written to be system independent.
 *	No checking for output buffer overflow is done, but flushes if needed!
 *	Returns nothing of value.
 */
func lprint(x int32) {
	debugf("(%d)", x)
	if len(lpbuf) >= cap(lpbuf) {
		lflush()
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(x))
	lpbuf = append(lpbuf, b[:]...)
}

/* macro to output one byte to the output buffer */
func lprc(ch byte) {
	if io_out != nil {
		lpbuf = append(lpbuf, ch)
		return
	}
	moved := false
	switch {
	case ch >= 32:
		win.Addch(ch)
		cursorX++
		moved = true
	case ch == '\n':
		// TODO: account for delete_line and insert_line
		cursorX, cursorY = 1, cursorY+1
		if cursorY >= 19 { // TODO: only if setscroll has been called?
			defer win.Clrtoeol()
		}
		moved = true
	case ch == '\t':
		n := 8 - (cursorX-1)%8
		win.Addstr(strings.Repeat(" ", n))
		cursorX += n
		moved = true
	}
	if moved {
		cursor(cursorX, cursorY)
	}
}

/* macro to turn on bold display for the terminal */
func setbold() {
	if boldon {
		win.AttrOn(curses.A_BOLD)
	} else {
		win.AttrOn(curses.A_REVERSE)
	}
}

/* macro to turn off bold display for the terminal */
func resetbold() {
	win.AttrOff(curses.A_BOLD) /* TODO: A_REVERSE too? */
}

/* macro to setup the scrolling region for the terminal */
func setscroll() {
	/* lprcat("\033[20;24r") */
	if err := win.SetScrollRegion(20, 24); err != nil {
		debugf("win.SetScrollRegion: %v", err)
	}
	if err := win.Scrollok(true); err != nil {
		debugf("win.Scrollok: %v", err)
	}
}

/* macro to clear the scrolling region for the terminal */
func resetscroll() {
	/* lprcat("\033[;24r") */
	if err := win.SetScrollRegion(0, 24); err != nil {
		debugf("win.SetScrollRegion: %v", err)
	}
	if err := win.Scrollok(false); err != nil {
		debugf("win.Scrollok: %v", err)
	}
}

/* macro to clear the screen and home the cursor */
func clear() {
	debugf("()")
	win.Clear()
	cursor(1, 1)
	win.Refresh()
	cbak[SPELLS] = -50
}

func cltoeoln() {
	win.Clrtoeol()
}

/*
 *	lwrite(buf,len)		write a buffer to the output buffer
 *		char *buf;
 *		int len;
 *
 *	Enter with the address and number of bytes to write out
 *	Returns nothing of value
 */
func lwrite(s string) {
	debugf("s=%q", s)
	if len(s) > 399 { /* don't copy data if can just write it */
		c[BYTESOUT] += len(s)

		//#ifndef VT100
		//		for (s := buf; len > 0; --len)
		//			lprc(*s++);
		//#else	/* VT100 */
		lflush()
		if _, err := io_out.WriteString(s); err != nil {
			debugf("Writing to output file %s: %v", io_out.Name(), err)
		}
		//#endif	/* VT100 */
	} else {
		for s != "" {
			if len(lpbuf) >= cap(lpbuf) {
				lflush() /* if buffer is full flush it	 */
			}
			// TODO: this isn't the correct computation, but it matches the original C version
			num2 := BUFBIG - len(lpbuf) /* # bytes left in output buffer	 */
			if num2 > len(s) {
				num2 = len(s)
			}
			lpbuf = append(lpbuf, []byte(s[:num2])...)
		}
	}
}

/*
 *	long lgetc()	Read one character from input buffer
 *
 *  Returns 0 if EOF, otherwise the character
 */
func lgetc() int {
	var buf [1]byte
	_, err := io_in.Read(buf[:])
	if err != nil {
		debugf("Reading from input file %s: %v", io_in.Name(), err)
		return 0
	}
	return int(buf[0])
}

/*
 *	long lrint()	Read one integer from input buffer
 *
 *		+---------+---------+---------+---------+
 *		|   high  |	    |	      |	  low	|
 *		|  order  |	    |	      |  order	|
 *		|   byte  |	    |	      |	  byte	|
 *		+---------+---------+---------+---------+
 *	       31  ---  24 23 --- 16 15 ---  8 7  ---   0
 *
 *	The save order is low order first, to high order (4 bytes total)
 *	Returns the int read
 */
func larn_lrint() int32 {
	var i uint32
	i = 255 & uint32(lgetc())
	i |= (255 & uint32(lgetc())) << 8
	i |= (255 & uint32(lgetc())) << 16
	i |= (255 & uint32(lgetc())) << 24
	return int32(i)
}

/*
 *	lrfill(address,number)		put input bytes into a buffer
 *		char *address;
 *		int number;
 *
 *	Reads "number" bytes into the buffer pointed to by "address".
 *	Returns nothing of value
 */
// TODO: implement lrfill?
/*
func lrfill(char *adr, int num) {
	u_char  *pnt;
	int    num2;

	while (num) {
		if (iepoint == ipoint) {
			if (num > 5) {	// fast way
				if (read(io_infd, adr, num) != num)
					write(2, "error reading from input file\n", 30);
				num = 0;
			} else {
				*adr++ = lgetc();
				--num;
			}
		} else {
			num2 = iepoint - ipoint;	// # of bytes left in the buffer
			if (num2 > num)
				num2 = num;
			pnt = inbuffer + ipoint;
			num -= num2;
			ipoint += num2;
			while (num2--)
				*adr++ = *pnt++;
		}
	}
}
*/

/*
 *	char *lgetw()			Get a whitespace ended word from input
 *
 *	Returns pointer to a buffer that contains word.  If EOF, returns a NULL
 */
func lgetw() (word string) {
	defer func() { debugf("=> %q", word) }()
	n, quote := LINBUFSIZE, 0
	var cc int
	lgp := ""
	for {
		cc = lgetc()
		if cc > 32 || cc == 0 {
			break
		}
	} /* eat whitespace */
	for {
		if cc == 0 && lgp == "" {
			return "" /* EOF */
		}
		if (n <= 1) || (cc <= 32 && quote == 0) {
			return lgp
		}
		if cc != '"' {
			lgp += string(cc)
		} else {
			quote ^= 1
		}
		n--
		cc = lgetc()
	}
	panic("unreachable")
}

/*
 *	char *lgetl()	Function to read in a line ended by newline or EOF
 *
 * Returns pointer to a buffer that contains the line.  If EOF, returns NULL
 */
func lgetl() (line string) {
	//defer func() { debugf("=> %q", line) }()
	i := LINBUFSIZE
	str := ""
	for {
		ch := lgetc()
		if ch != 0 {
			str += string(ch)
		}
		if ch == 0 {
			if str == "" {
				return "" /* EOF */
			}
			return str /* line ended by EOF */
		}
		if ch == '\n' || i <= 1 {
			return str /* line ended by \n */
		}
		i--
	}
	panic("unreachable")
}

/*
 *	lcreat(filename)			Create a new file for write
 *
 *	lcreat("") means to the terminal
 *	Returns false if error.
 */
func lcreat(str string) bool {
	debugf("%q", str)
	lflush()
	// TODO: original C version shrinks the output buffer to BUFBIG, despite allocating more. why?
	lpbuf = lpbuf[:0]
	if str == "" {
		io_out = nil
		return true
	}
	var err error
	io_out, err = os.OpenFile(str, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		io_out = nil
		debugf("Creating file %s: %v", str, err)
		lflush() // TODO: needed?
		return false
	}
	return true
}

/*
 *	lopen(filename)			Open a file for read
 *
 *	lopen("") means from the terminal
 *	Returns false if error.
 */
func lopen(str string) bool {
	debugf("%q", str)
	if str == "" {
		io_in = os.Stdin
		return true
	}
	var err error
	io_in, err = os.Open(str)
	if err != nil {
		lwclose()
		io_out = nil
		lpbuf = lpbuf[:0]
		return false
	}
	return true
}

/*
 *	lappend(filename)		Open for append to an existing file
 *
 *	lappend("") means to the terminal
 *	Returns false if error.
 */
func lappend(str string) bool {
	debugf("%q", str)
	if str == "" {
		io_out = nil
		return true
	}
	var err error
	io_out, err = os.OpenFile(str, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		debugf("Opening for append output file %s: %v", str, err)
		io_out = nil
		return false
	}
	//lseek(io_outfd, 0, SEEK_END);	/* seek to end of file */
	return true
}

/*
 *	lrclose() close the input file
 *
 *	Returns nothing of value.
 */
func lrclose() {
	debugf("()")
	if io_in != os.Stdin {
		if err := io_in.Close(); err != nil {
			debugf("Closing input file %s: %v", io_in.Name(), err)
		}
		io_in = os.Stdin
	}
}

/*
 *	lwclose() close output file flushing if needed
 *
 *	Returns nothing of value.
 */
func lwclose() {
	debugf("()")
	lflush()
	if io_out != nil && io_out != os.Stderr {
		if err := io_out.Close(); err != nil {
			debugf("Closing output file %s: %v", io_out.Name(), err)
		}
		io_out = nil
	}
}

/*
 *	lprcat(string)	append a string to the output buffer
 *			    	avoids calls to lprintf (time consuming)
 */
func lprcat(str string) {
	debugf("(%q) @ (%d, %d)", str, cursorX, cursorY)
	// TODO: implement lprcat/lprc more efficiently?

	// Lazy way to support certain terminal escape codes.
	for {
		i := strings.Index(str, "\033")
		if i < 0 {
			break
		}
		lprcat(str[:i])
		str = str[i:]
		any := false
		for seq, f := range termCodes {
			if strings.HasPrefix(str[1:], seq) {
				f()
				any = true
				str = str[1+len(seq):]
				break
			}
		}
		if !any {
			break
		}
	}

	for _, ch := range str {
		lprc(byte(ch))
	}
}

var termCodes = map[string]func(){
	"[m":  func() { win.AttrOff(curses.A_BOLD | curses.A_REVERSE | curses.A_STANDOUT | curses.A_UNDERLINE) },
	"[4m": func() { win.AttrOn(curses.A_UNDERLINE) },
	"[7m": func() { win.AttrOn(curses.A_REVERSE) },
}

/*
 *	cursor(x,y) 		Subroutine to set the cursor position
 *
 *	x and y are the cursor coordinates, and lpbuff is the output buffer where
 *	escape sequence will be placed.
 */

var cursorX, cursorY int

func cursor(x, y int) {
	//debugf("(%d, %d)", x, y)
	// curses is based on [0,0]. The rest of Larn is based on [1,1]
	win.Move(x-1, y-1)
	cursorX, cursorY = x, y
}

/*
 *	Routine to position cursor at beginning of 24th line
 */
func cursors() {
	cursor(1, 24)
}

//#ifndef VT100
/*
 * Warning: ringing the bell is control code 7. Don't use in defines.
 * Don't change the order of these defines.
 * Also used in helpfiles. Codes used in helpfiles should be \E[1 to \E[7 with
 * obvious meanings.
 */

//static char    *outbuf = 0;     /* translated output buffer */
/*
 * init_term()		Terminal initialization -- setup termcap info
 */
func init_term() {
	debugf("()")
	w, err := curses.Initscr()
	if err != nil {
		log.Fatalf("curses.Initscr: %v", err)
	}
	win = w
	/*
		if err := curses.Noecho(); err != nil {
			log.Fatalf("curses.Noecho: %v", err)
		}
	*/
	/*
		setupterm(NULL, 0, NULL); // will exit if invalid term
		if (!cursor_address) {
			fprintf(stderr, "term does not have cursor_address.\n");
			exit(1);
		}
		if (!clr_eol) {
			fprintf(stderr, "term does not have clr_eol.\n");
			exit(1);
		}
		if (!clear_screen) {
			fprintf(stderr, "term does not have clear_screen.\n");
			exit(1);
		}
		if ((outbuf = malloc(BUFBIG + 16)) == 0) {      // get memory for decoded output buffer
		    fprintf(stderr, "Error malloc'ing memory for decoded output buffer\n");
		    died(-285);     // malloc() failure
		}
	*/
}

/*
 * cl_line(x,y)  Clear the whole line indicated by 'y' and leave cursor at [x,y]
 */
func cl_line(x, y int) {
	debugf("(%d, %d)", x, y)
	cursor(1, y)
	win.Clrtoeol()
	cursor(x, y)
}

/*
 * cl_up(x,y) Clear screen from [x,1] to current position. Leave cursor at [x,y]
 */
func cl_up(x, y int) {
	debugf("(%d, %d)", x, y)
	for i := 1; i <= y; i++ {
		cursor(1, i)
		win.Clrtoeol()
	}
	cursor(x, y)
}

/*
 * cl_dn(x,y) 	Clear screen from [1,y] to end of display. Leave cursor at [x,y]
 */
func cl_dn(x, y int) {
	debugf("(%d, %d)", x, y)
	// TODO: factor in clr_eos
	for i := y; i <= 24; i++ {
		cursor(1, i)
		win.Clrtoeol()
	}
	cursor(x, y)
}

/*
 * standout(str)	Print the argument string in inverse video (standout mode).
 */
func standout(str string) {
	debugf("(%q)", str)
	win.AttrOn(curses.A_STANDOUT) // TODO: or A_REVERSE?
	lprcat(str)
	win.AttrOff(curses.A_STANDOUT) // TODO: or A_REVERSE?
}

/*
 * set_score_output() 	Called when output should be literally printed.
 */
func set_score_output() {
	enable_scroll = -1
}

/*
 *	lflush()	Flush the output buffer
 *
 *	Returns nothing of value.
 *	for termcap version: Flush output in output buffer according to output
 *	status as indicated by `enable_scroll'
 */
//#ifndef VT100
var scrline = 18 /* line # for wraparound instead of scrolling if no DL */
func lflush() {
	debugf("()")
	if win != nil { // set to nil when exiting
		win.Refresh()
	}
	/*
		int    lpoint;
		u_char  *str;
		static int      curx = 0;
		static int      cury = 0;

		if ((lpoint = lpnt - lpbuf) > 0) {
			c[BYTESOUT] += lpoint;

			if (enable_scroll <= -1) {
				flush_buf();
				if (write(io_outfd, lpbuf, lpoint) != lpoint)
					write(2, "error writing to output file\n", 29);
				lpnt = lpbuf;	// point back to beginning of buffer
				return;
			}
			for (str = lpbuf; str < lpnt; str++) {
				if (*str >= 32) {
					ttputch(*str);
					curx++;
				} else
					switch (*str) {
					case CLEAR:
						tputs(clear_screen, 0, ttputch);
						curx = cury = 0;
						break;

					case CL_LINE:
						tputs(clr_eol, 0, ttputch);
						break;

					case CL_DOWN:
						tputs(clr_eos, 0, ttputch);
						break;

					case ST_START:
						tputs(enter_standout_mode, 0, ttputch);
						break;

					case ST_END:
						tputs(exit_standout_mode, 0, ttputch);
						break;

					case CURSOR:
						curx = *++str - 1;
						cury = *++str - 1;
						tputs(tiparm(cursor_address,
							    cury, curx), 0, ttputch);
						break;

					case '\n':
						if ((cury == 23) && enable_scroll) {
							if (!delete_line ||
							    !insert_line)
							{	// wraparound or scroll?
								if (++scrline > 23)
									scrline = 19;

								if (++scrline > 23)
									scrline = 19;
								tputs(tiparm(
								    cursor_address,
								    scrline, 0),
								    0, ttputch);
								tputs(clr_eol, 0,
								    ttputch);

								if (--scrline < 19)
									scrline = 23;
								tputs(tiparm(
								    cursor_address,
								    scrline, 0),
								    0, ttputch);
								tputs(clr_eol, 0,
								    ttputch);
							} else {
								tputs(tiparm(
								    cursor_address,
								    19, 0),
								    0, ttputch);
								tputs(delete_line, 0,
								    ttputch);
								tputs(tiparm(
								    cursor_address,
								    23, 0),
								    0, ttputch);
								//
								// tputs (AL, 0,
								// ttputch);
								//
							}
						} else {
							ttputch('\n');
							cury++;
						}
						curx = 0;
						break;

					default:
						ttputch(*str);
						curx++;
					};
			}
		}
		lpnt = lpbuf;
		flush_buf();		// flush real output buffer now
	*/
}

//#else	// VT100 */
/*
void
lflush()
{
	int    lpoint;
	if ((lpoint = lpnt - lpbuf) > 0) {
		c[BYTESOUT] += lpoint;

		if (write(io_outfd, lpbuf, lpoint) != lpoint)
			write(2, "error writing to output file\n", 29);
	}
	lpnt = lpbuf;		// point back to beginning of buffer
}
//#endif	// VT100
*/

//#ifndef VT100
var vindex = 0

/*
 * ttputch(ch)		Print one character in decoded output buffer.
 */
func ttputch(ch int) int {
	debugf("(%c)", ch)
	// TODO: implement ttputch?
	/*
		outbuf[vindex++] = ch;
		if (vindex >= BUFBIG)
			flush_buf();
	*/
	return 0
}

/*
 * flush_buf()			Flush buffer with decoded output.
 */
func flush_buf() {
	//if (vindex)
	//	write(io_outfd, outbuf, vindex);
	//vindex = 0;
}

/*
 *	beep()	Routine to emit a beep if enabled (see no-beep in .larnopts)
 */
func beep() {
	if len(replayActions) > 0 {
		return
	}
	if !nobeep {
		curses.Beep()
	}
}
