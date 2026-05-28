import { me as fetchMe } from '$lib/api/auth-api';
import type { AuthUser, SiteRole } from '$lib/types/auth';

function siteRole(user: AuthUser | null): SiteRole | undefined {
	if (!user) return undefined;
	return user.siteRole ?? user.role;
}

class Auth {
	#user = $state<AuthUser | null>(null);
	#loading = $state(true);
	#error = $state<string | null>(null);

	get user() {
		return this.#user;
	}

	get loading() {
		return this.#loading;
	}

	get error() {
		return this.#error;
	}

	get isAdmin() {
		const r = siteRole(this.#user);
		return r === 'site_admin';
	}

	async load(): Promise<void> {
		this.#loading = true;
		this.#error = null;

		try {
			this.#user = await fetchMe();
		} catch (err) {
			this.#user = null;
			this.#error = err instanceof Error ? err.message : 'Failed to load user';
		} finally {
			this.#loading = false;
		}
	}

	setUser(user: AuthUser | null): void {
		this.#user = user;
		this.#error = null;
		this.#loading = false;
	}

	clear(): void {
		this.#user = null;
		this.#error = null;
		this.#loading = false;
	}
}

export const auth = new Auth();
