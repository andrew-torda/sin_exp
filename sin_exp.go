// Program sin_exp generates data for the R exercise. It takes a
// Matrikelnummer as input for the random see and uses it to set
// phase, frequency, and decay rate in a sin (2 pi x + phi) * exp (-bx)
// function.
// If the secretstring is set in the environment, the chosen values will
// be printed to stdout. This is for correcting the exercise.
// -d option tells us not to use the matrikelnumm from the command line. Use a default value.
// -n do not write the output file. Useful if you just want the parameters corresponding to a matrikelnummer.
// -o fname, write the xy output to fname.
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"sort"
	"strconv"
)

const (
	exitSuccess = 0
	exitFailure = iota
)

const (
	freqIni       = 100      // omega, frequency
	decayIni      = 30       // decay rate
	randVariation = 0.10     // amount of random variation to add
	xmax          = 0.1      // x values from 0 to xmax
	npoint        = 500      // Generate this many points
	noisefrac     = 0.05     // Add this much noise as a fraction of maximum range
	secretstring  = "qwerty" // Look for this string in environment, print out params
)

type xypair struct {
	x float32
	y float32
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:", path.Base(os.Args[0]), "matrikelnummer")
	flag.PrintDefaults()
}

// add_variation takes the first number and adds or substracts
// a random fraction, given by the second argument
func addVariation(ini, vary float32) float32 {
	r := vary * (2*rand.Float32() - 1)
	return ini + r*ini
}

// randInRange returns a quasi-random number within a given range
func randInRange(min, max float32) float32 {
	r := (max - min) * rand.Float32()
	return min + r
}

// getFuncParams returns the phase, freq and decay rate based on
// the xxxIni constants + random numbers
func getFuncParams() (phase, freq, decay float32) {
	phase = randInRange(0, 2*math.Pi)
	freq = addVariation(freqIni, randVariation)
	decay = addVariation(decayIni, randVariation)
	return
}

// makeSinExp takes the phase shift, frequency and decay rate and
// returns a set of xypairs from sin(x) * exp (x * decay)
func makeSinExp(xypairs []xypair, phase, freq, decay float32) {
	twoPi := 2 * math.Pi
	phase64 := float64(phase)
	for i := range xypairs {
		x := xypairs[i].x
		xp := twoPi * float64(x*freq)
		y := math.Sin(xp + phase64)
		y = y * math.Exp(-float64(decay*x))
		xypairs[i].y = float32(y)
	}
}

// fillX generates x-values for the pairs. These are random numbers within
// an interval. We sort them since it lets us use lines when plotting
// in the Uebung. Could also have been done using a poisson distribution
func fillX(xypairs []xypair, xmax float32) {
	for i := range xypairs {
		xypairs[i].x = randInRange(0, xmax)
	}
	sort.Slice(xypairs, func(i, j int) bool {
		return xypairs[i].x > xypairs[j].x
	})
}

// noise adds noise to the y values. It is calculated as a fraction
// of the range of y values.
func noise(xypairs []xypair, fracnoise float32) {
	minval := xypairs[0].y
	maxval := minval
	for _, xy := range xypairs {
		if xy.y < minval {
			minval = xy.y
		}
		if xy.y > maxval {
			maxval = xy.y
		}
	}
	minval = minval * fracnoise
	maxval = maxval * fracnoise
	for i := range xypairs {
		xypairs[i].y = xypairs[i].y + randInRange(minval, maxval)
	}
}

// writeXy writes the xy pairs to outfile or standard output if outfile is
// not defined. If the nowriteFlag is set, we just return.
func writeXy(outfile *string, xypairs []xypair, nowriteFlag *bool) error {
	if *nowriteFlag {
		return nil
	}
	var fp *os.File
	if *outfile != "" {
		fmt.Println("Writing results to ", *outfile)
		_, err := os.Stat(*outfile)
		if err == nil {
			fmt.Fprintln(os.Stderr, "File", *outfile, "exists... Overwriting it.")
		}
		fp, err = os.Create(*outfile)
		if err != nil {
			return err
		}
		defer fp.Close()
	} else {
		fp = os.Stdout
	}
	fmt.Fprintln(fp, "x y")
	for _, xy := range xypairs {
		fmt.Fprintln(fp, xy.x, xy.y)
	}
	return nil
}

func mymain() int {
	outfile := flag.String("o", "", "file for table output")
	dflt_seed := flag.Bool("d", false, "use default rather than Matrikelnummer")
	nowriteFlag := flag.Bool("n", false, "do not write the xyfile")
	flag.Parse()
	if len(flag.Args()) < 1 && *dflt_seed == false {
		usage()
		return (exitFailure)
	}

	{
		var seed int64
		if *dflt_seed {
			seed = 1
		} else {
			s := flag.Arg(0)
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Println(err)
				return (exitFailure)
			}
			seed = i
			if seed > 9999999 || seed < 111111 {
				fmt.Fprintln(os.Stderr, s, "is probably not a valid matrikelnummer")
				return (exitFailure)
			}
		}
		rand.Seed(seed)
	}
	phase, freq, decay := getFuncParams()
	_, printValFlag := os.LookupEnv(secretstring)
	if printValFlag {
		fmt.Println("phase", phase, "freq", freq, "decay", decay)
	}

	xypairs := make([]xypair, npoint)
	fillX(xypairs, xmax)
	makeSinExp(xypairs, phase, freq, decay)
	noise(xypairs, noisefrac)
	err := writeXy(outfile, xypairs, nowriteFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed writing to", *outfile, err)
		return (exitFailure)
	}
	return (exitSuccess)
}

func main() {
	os.Exit(mymain())
}
