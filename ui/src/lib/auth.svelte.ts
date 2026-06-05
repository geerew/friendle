import { me as fetchMe } from '$lib/api/auth-api';
import type { AuthUser, SiteRole } from '$lib/types/auth';

function siteRole(user: AuthUser | null): SiteRole | undefined {
	if (!user) return undefined;
	return user.siteRole ?? user.role;
}

class Auth {
	#user = $state<AuthUser | null>(null);
	#loading = $state(false);
	#initialized = $state(false);
	#error = $state<string | null>(null);

	get user() {
		return this.#user;
	}

	get loading() {
		return this.#loading;
	}

	get initialized() {
		return this.#initialized;
	}

	get error() {
		return this.#error;
	}

	get isAdmin() {
		const r = siteRole(this.#user);
		return r === 'site_admin';
	}

	// load fetches the current user once on cold start; later updates come from setUser/clear
	async load(): Promise<void> {
		if (this.#initialized) {
			return;
		}

		this.#loading = true;
		this.#error = null;

		try {
			await new Promise((resolve) => setTimeout(resolve, 500));
			this.#user = await fetchMe();
		} catch (err) {
			this.#user = null;
			this.#error = err instanceof Error ? err.message : 'Failed to load user';
		} finally {
			this.#loading = false;
			this.#initialized = true;
		}
	}

	setUser(user: AuthUser | null): void {
		this.#user = user;
		this.#error = null;
		this.#loading = false;
		this.#initialized = true;
	}

	clear(): void {
		this.#user = null;
		this.#error = null;
		this.#loading = false;
		this.#initialized = true;
	}
}

export const auth = new Auth();
