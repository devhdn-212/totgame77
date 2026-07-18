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

// Per-type max bet and DISC discount rate for the base 4D/3D/2D bet types.
export const BET_TYPE_LIMITS: Record<string, { maxBet: number; discRate: number }> = {
	"4D": { maxBet: 50000, discRate: 0.66 },
	"3D": { maxBet: 200000, discRate: 0.59 },
	"3DD": { maxBet: 200000, discRate: 0.56 },
	"2D": { maxBet: 1000000, discRate: 0.29 },
	"2DD": { maxBet: 1000000, discRate: 0.26 },
	"2DT": { maxBet: 1000000, discRate: 0.26 },
};

// DISC: payout = bet - (bet * discRate). FULL/BB: payout = bet, no discount.
export function calculatePayout(
	type: string,
	bet: Decimal.Value,
	kombinasi?: string,
): { disc: Decimal | null; payout: Decimal } {
	const betDecimal = new Decimal(bet);
	const limits = BET_TYPE_LIMITS[type];
	if (kombinasi === "DISC" && limits) {
		const disc = betDecimal.times(limits.discRate);
		return { disc, payout: betDecimal.minus(disc) };
	}
	return { disc: null, payout: betDecimal };
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, "children"> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };
