package main

import "math"

// This is a port of Ken Perlin's Java code. The
// original Java code is at http://cs.nyu.edu/%7Eperlin/noise/.
// Note that in this version, a number from 0 to 1 is returned.

func perlinNoise(x, y, z float64) float64 {
    p := make([]int, 512)

    permutation := [256]int{151, 160, 137, 91, 90, 15,
        131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
        190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
        88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
        77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
        102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
        135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
        5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
        223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
        129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
        251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
        49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
        138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
    }

    for i := range 256 {
        p[i] = permutation[i]
        p[256 + i] = permutation[i]
    }

    var (
    	X = int(math.Floor(x)) & 255                  // FIND UNIT CUBE THAT
        Y = int(math.Floor(y)) & 255                  // CONTAINS POINT.
        Z = int(math.Floor(z)) & 255
    )
    x -= math.Floor(x);                                // FIND RELATIVE X,Y,Z
    y -= math.Floor(y);                                // OF POINT IN CUBE.
    z -= math.Floor(z);

    var (
    	u = fade(x)                                // COMPUTE FADE CURVES
        v = fade(y)                                // FOR EACH OF X,Y,Z.
        w = fade(z)
    )
    var (
    	A = p[X] + Y
     	AA = p[A] + Z
     	AB = p[A + 1] + Z      // HASH COORDINATES OF
        B = p[X + 1] + Y
        BA = p[B] + Z
        BB = p[B + 1] + Z      // THE 8 CUBE CORNERS,
    )

    return scale(lerp(w, lerp(v, lerp(u, grad(p[AA], x, y, z),  // AND ADD
        grad(p[BA], x - 1, y, z)), // BLENDED
        lerp(u, grad(p[AB], x, y - 1, z),  // RESULTS
            grad(p[BB], x - 1, y - 1, z))),// FROM  8
        lerp(v, lerp(u, grad(p[AA + 1], x, y, z - 1),  // CORNERS
            grad(p[BA + 1], x - 1, y, z - 1)), // OF CUBE
            lerp(u, grad(p[AB + 1], x, y - 1, z - 1),
                grad(p[BB + 1], x - 1, y - 1, z - 1)))));
}

func fade(t float64) float64 {
    return t * t * t * (t * (t * 6 - 15) + 10);
}

func lerp(t, a, b float64) float64 {
    return a + t * (b - a);
}

func grad(hash int, x, y, z float64) float64 {
    // CONVERT LO 4 BITS OF HASH CODE
    // INTO 12 GRADIENT DIRECTIONS.
    h := hash & 15

    var u, v float64
    if h < 8 {
    	u = x
    } else {
     	u = y
    }

    if h < 4 {
    	v = y
    } else if h == 12 || h == 14 {
    	v = x
    } else {
    	v = z
    }

    term1 := -u
    if (h & 1) == 0 {
    	term1 = u
    }

    term2 := -v
    if (h & 2) == 0 {
    	term2 = v
    }

    return term1 + term2
}

func scale(n float64) float64 {
    return (1 + n) / 2;
}
