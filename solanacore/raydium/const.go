package raydium

const (
	LiquidityPoolProgramV2  = "RVKd61ztZW9GUwhRbbLoYVRE5Xf1B2tVscKqwZqXgEr"
	LiquidityPoolProgramV3  = "27haf8L6oxUeXrHrgEgsexjSY5hbVUWEmvv9Nyxg8vQv"
	LiquidityPoolProgramV4  = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
	LiquidityPoolProgramAMM = "5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h"
	ProgramAMMRouting       = "routeUGWgWzqBWFcrCfv8tritsqukccJPu3q5GPP3xS"
)

type SwapDirection uint8

const (
	DirectionUnknown SwapDirection = 0
	Pc2Coin          SwapDirection = 1
	Coin2Pc          SwapDirection = 2
)

const (
	LIQUIDITY_FEES_NUMERATOR   = 25
	LIQUIDITY_FEES_DENOMINATOR = 10000
)
