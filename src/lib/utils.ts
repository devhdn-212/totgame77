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

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, "children"> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };
