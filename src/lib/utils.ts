import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import Decimal from "decimal.js";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function formatIDR(value: Decimal.Value): string {
	const amount = new Decimal(value).toDecimalPlaces(0, Decimal.ROUND_HALF_UP);
	const digits = amount.abs().toFixed(0);
	const grouped = digits.replace(/\B(?=(\d{3})+(?!\d))/g, ".");
	return (amount.isNegative() ? "-Rp " : "Rp ") + grouped;
}

type BetRates = { discRate: number; keiRate?: number };

// Per-type max bet, DISC discount rate, and (for 5050/Dasar) the "kei"
// adjustment rate. keiRate is negative by design — it ADDS to payout when
// subtracted. perNumber overrides the base rates for specific selections
// (Dasar's Besar/Kecil/Genap/Ganjil each have their own kei rate).
// keiNegativeStyle flips the Kei line's display to red/parens (Dasar) instead
// of the default blue (5050).
export const BET_TYPE_LIMITS: Record<
	string,
	BetRates & { maxBet: number; perNumber?: Record<string, BetRates>; keiNegativeStyle?: boolean }
> = {
	"4D": { maxBet: 50000, discRate: 0.66 },
	"3D": { maxBet: 200000, discRate: 0.59 },
	"3DD": { maxBet: 200000, discRate: 0.56 },
	"2D": { maxBet: 1000000, discRate: 0.29 },
	"2DD": { maxBet: 1000000, discRate: 0.26 },
	"2DT": { maxBet: 1000000, discRate: 0.26 },
	COLOK_BEBAS: { maxBet: 3000000, discRate: 0.06 },
	COLOK_MACAU: { maxBet: 3000000, discRate: 0.1 },
	COLOK_NAGA: { maxBet: 3000000, discRate: 0.1 },
	COLOK_JITU: { maxBet: 3000000, discRate: 0.05 },
	"50_50_UMUM": { maxBet: 1000000, discRate: 0.02, keiRate: -0.015 },
	"50_50_SPECIAL": { maxBet: 1000000, discRate: 0.015, keiRate: -0.025 },
	"50_50_KOMBINASI": { maxBet: 1000000, discRate: 0.015, keiRate: -0.025 },
	MACAU_KOMBINASI: { maxBet: 1000000, discRate: 0.05 },
	SHIO: { maxBet: 1000000, discRate: 0.08 },
	DASAR: {
		maxBet: 1000000,
		discRate: 0,
		keiNegativeStyle: true,
		perNumber: {
			BESAR: { discRate: 0, keiRate: -0.25 },
			KECIL: { discRate: 0, keiRate: 0 },
			GENAP: { discRate: 0, keiRate: -0.25 },
			GANJIL: { discRate: 0, keiRate: 0 },
		},
	},
};

// DISC: payout = bet - (bet * discRate). FULL/BB: payout = bet, no discount.
// Types with a keiRate additionally apply: payout = bet - (bet * keiRate) - (bet * discRate).
export function calculatePayout(
	type: string,
	number: string,
	bet: Decimal.Value,
	kombinasi?: string,
): { kei: Decimal | null; disc: Decimal | null; payout: Decimal; keiNegativeStyle: boolean } {
	const betDecimal = new Decimal(bet);
	const limits = BET_TYPE_LIMITS[type];
	if (kombinasi === "DISC" && limits) {
		const rates = limits.perNumber?.[number] ?? limits;
		const disc = betDecimal.times(rates.discRate);
		if (rates.keiRate !== undefined) {
			const kei = betDecimal.times(rates.keiRate).negated();
			return {
				kei,
				disc,
				payout: betDecimal.plus(kei).minus(disc),
				keiNegativeStyle: !!limits.keiNegativeStyle,
			};
		}
		return { kei: null, disc, payout: betDecimal.minus(disc), keiNegativeStyle: false };
	}
	return { kei: null, disc: null, payout: betDecimal, keiNegativeStyle: false };
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, "children"> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };
