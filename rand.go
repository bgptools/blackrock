package blackrock

import "math"

var sbox = [256]byte{
	0x91, 0x58, 0xb3, 0x31, 0x6c, 0x33, 0xda, 0x88,
	0x57, 0xdd, 0x8c, 0xf2, 0x29, 0x5a, 0x08, 0x9f,
	0x49, 0x34, 0xce, 0x99, 0x9e, 0xbf, 0x0f, 0x81,
	0xd4, 0x2f, 0x92, 0x3f, 0x95, 0xf5, 0x23, 0x00,
	0x0d, 0x3e, 0xa8, 0x90, 0x98, 0xdd, 0x20, 0x00,
	0x03, 0x69, 0x0a, 0xca, 0xba, 0x12, 0x08, 0x41,
	0x6e, 0xb9, 0x86, 0xe4, 0x50, 0xf0, 0x84, 0xe2,
	0xb3, 0xb3, 0xc8, 0xb5, 0xb2, 0x2d, 0x18, 0x70,

	0x0a, 0xd7, 0x92, 0x90, 0x9e, 0x1e, 0x0c, 0x1f,
	0x08, 0xe8, 0x06, 0xfd, 0x85, 0x2f, 0xaa, 0x5d,
	0xcf, 0xf9, 0xe3, 0x55, 0xb9, 0xfe, 0xa6, 0x7f,
	0x44, 0x3b, 0x4a, 0x4f, 0xc9, 0x2f, 0xd2, 0xd3,
	0x8e, 0xdc, 0xae, 0xba, 0x4f, 0x02, 0xb4, 0x76,
	0xba, 0x64, 0x2d, 0x07, 0x9e, 0x08, 0xec, 0xbd,
	0x52, 0x29, 0x07, 0xbb, 0x9f, 0xb5, 0x58, 0x6f,
	0x07, 0x55, 0xb0, 0x34, 0x74, 0x9f, 0x05, 0xb2,

	0xdf, 0xa9, 0xc6, 0x2a, 0xa3, 0x5d, 0xff, 0x10,
	0x40, 0xb3, 0xb7, 0xb4, 0x63, 0x6e, 0xf4, 0x3e,
	0xee, 0xf6, 0x49, 0x52, 0xe3, 0x11, 0xb3, 0xf1,
	0xfb, 0x60, 0x48, 0xa1, 0xa4, 0x19, 0x7a, 0x2e,
	0x90, 0x28, 0x90, 0x8d, 0x5e, 0x8c, 0x8c, 0xc4,
	0xf2, 0x4a, 0xf6, 0xb2, 0x19, 0x83, 0xea, 0xed,
	0x6d, 0xba, 0xfe, 0xd8, 0xb6, 0xa3, 0x5a, 0xb4,
	0x48, 0xfa, 0xbe, 0x5c, 0x69, 0xac, 0x3c, 0x8f,

	0x63, 0xaf, 0xa4, 0x42, 0x25, 0x50, 0xab, 0x65,
	0x80, 0x65, 0xb9, 0xfb, 0xc7, 0xf2, 0x2d, 0x5c,
	0xe3, 0x4c, 0xa4, 0xa6, 0x8e, 0x07, 0x9c, 0xeb,
	0x41, 0x93, 0x65, 0x44, 0x4a, 0x86, 0xc1, 0xf6,
	0x2c, 0x97, 0xfd, 0xf4, 0x6c, 0xdc, 0xe1, 0xe0,
	0x28, 0xd9, 0x89, 0x7b, 0x09, 0xe2, 0xa0, 0x38,
	0x74, 0x4a, 0xa6, 0x5e, 0xd2, 0xe2, 0x4d, 0xf3,
	0xf4, 0xc6, 0xbc, 0xa2, 0x51, 0x58, 0xe8, 0xae,
}

type Blackrock struct {
	cyclerange int
	a          int
	b          int
	seed       int
	rounds     int
	aBits      int
	aMask      int
	bBits      int
	bMask      int
}

func Init(cyclerange int, seed int, rounds int) *Blackrock {
	br := Blackrock{}
	foo := math.Sqrt(float64(cyclerange) * 1.0)

	/* This algorithm gets very non-random at small numbers, so I'm going
	 * to try to fix some constants here to make it work. It doesn't have
	 * to be good, since it's kinda pointless having ranges this small */
	switch cyclerange {
	case 0:
		br.a = 0
		br.b = 0
		break
	case 1:
		br.a = 1
		br.b = 1
		break
	case 2:
		br.a = 1
		br.b = 2
		break
	case 3:
		br.a = 2
		br.b = 2
		break
	case 4:
	case 5:
	case 6:
		br.a = 2
		br.b = 3
		break
	case 7:
	case 8:
		br.a = 3
		br.b = 3
		break
	default:
		br.cyclerange = cyclerange
		br.a = int(foo - 2)
		br.b = int(foo + 3)
		break
	}

	for br.a*br.b <= cyclerange {
		br.b++
	}

	br.rounds = rounds
	br.seed = seed
	br.cyclerange = cyclerange
	return &br
}

func encrypt(r int, a int, b int, m int, seed uint) int {
	var L, R, tmp int

	L = m % a
	R = m / a

	for j := 1; j <= r; j++ {
		if j&1 != 0 {
			tmp = (L + read(j, R, seed)) % a
		} else {
			tmp = (L + read(j, R, seed)) % b
		}
		L = R
		R = tmp
	}
	if r&1 != 0 {
		return a*L + R
	} else {
		return a*R + L
	}
}

func getbyte(R, n int, seed uint, r uint) uint {
	RR := uint(R)
	nn := uint(n)

	return (((RR) >> (nn * 8)) ^ seed ^ r) & 0xFF
}

func read(r, R int, seed uint) int {
	var r0, r1, r2, r3 int
	rr := uint(r)

	R ^= int((seed << rr) ^ (seed >> (64 - rr)))

	r0 = int(sbox[getbyte(R, 0, seed, uint(r))]<<0 | sbox[getbyte(R, 1, seed, uint(r))]<<8)
	r1 = int((sbox[getbyte(R, 2, seed, uint(r))]<<16 | sbox[getbyte(R, 3, seed, uint(r))]<<24))
	r2 = int(sbox[getbyte(R, 4, seed, uint(r))]<<0 | sbox[getbyte(R, 5, seed, uint(r))]<<8)
	r3 = int((sbox[getbyte(R, 6, seed, uint(r))]<<16 | sbox[getbyte(R, 7, seed, uint(r))]<<24))

	R = r0 ^ r1 ^ r2<<23 ^ r3<<33

	return R
}

func unencrypt(r uint, a, b, m, seed int) int {

	var L, R int
	var j uint
	var tmp int

	if r&1 != 0 {
		R = m % a
		L = m / a
	} else {
		L = m % a
		R = m / a
	}

	for j = r; j >= 1; j-- {
		if j&1 != 0 {
			tmp = read(int(j), int(L), uint(seed))
			if tmp > R {
				tmp = (tmp - R)
				tmp = a - (tmp % a)
				if tmp == a {
					tmp = 0

				}
			} else {
				tmp = (R - tmp)
				tmp %= a
			}
		} else {
			tmp = read(int(j), int(L), uint(seed))
			if tmp > R {
				tmp = (tmp - R)
				tmp = b - (tmp % b)
				if tmp == b {
					tmp = 0
				}
			} else {
				tmp = (R - tmp)
				tmp %= b
			}
		}
		R = L
		L = tmp
	}
	return a*R + L
}

func (br *Blackrock) Shuffle(m int) int {
	var c int

	c = encrypt(br.rounds, br.a, br.b, m, uint(br.seed))
	for c >= br.cyclerange {
		c = encrypt(br.rounds, br.a, br.b, c, uint(br.seed))

	}

	return c
}
