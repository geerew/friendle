export type SiteRole = 'site_admin' | 'site_user';

export interface SignupStatus {
	enabled: boolean;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	displayName?: string;
	password: string;
}

export interface SelfUpdateRequest {
	displayName?: string;
	currentPassword?: string;
	password?: string;
}

export interface SelfDeleteRequest {
	currentPassword: string;
}

export interface AuthUser {
	id: string;
	username: string;
	displayName: string;
	siteRole: SiteRole;
	/** @deprecated use siteRole */
	role?: SiteRole;
}

export type User = AuthUser;
