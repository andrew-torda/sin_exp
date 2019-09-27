# 28 June 2018

# Input data
datafile='./sinexp.dat'
outputplotfile='uebung_2.png'

sinexp <- function (x, phase, freq, decay) {
     v <- sin(2 * pi * x * freq + phase)
     v <- v * exp(-decay * x)
}

# Put phase back in range of 0 to 2 pi
squash_to_range <- function (x) {
    while (x > 2*pi) { x = x - 2 * pi}
    while (x < 0)    { x = x + 2 * pi}
    return (x)
}
# main

d <- read.table (datafile, header=TRUE)

nlmod <- nls (y~sinexp(x, phase, freq, decay), data = d,
              start = list(phase=4.1, freq=150, decay=30))
#png(file=outputplotfile)
plot (d$x, d$y, xlab = expression (italic(x)), ylab="amplitude")
pred=predict (nlmod)
lines(d$x, predict(nlmod), col = 3, lwd = 3)
cc <- coef(nlmod)
phase <- squash_to_range(cc['phase']) # bring phase back within 0 to 2 pi
fancy <- paste( "freq", sprintf("%.1f", cc['freq']),
               "\nphase", sprintf("%.1f", phase),
               "\ndecay ", sprintf("%.1f", cc['decay']))

text (x=0.05, y=0.5, labels=c(fancy), adj = 0)
#dev.off()
nlmod
