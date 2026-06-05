import { apiFetch, parseJson } from './fetch';
import type {
	LoginRequest,
	RegisterRequest,
	SelfDeleteRequest,
	SelfUpdateRequest,
	SignupStatus,
	User
} from '$lib/types/auth';

export async function getSignupStatus(): Promise<SignupStatus> {
	const response = await apiFetch('/api/auth/signup-status');
	return parseJson(response);
}

export async function login(data: LoginRequest): Promise<User> {
	const response = await apiFetch('/api/auth/login', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function register(data: RegisterRequest): Promise<User> {
	const response = await apiFetch('/api/auth/register', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function bootstrap(token: string, data: RegisterRequest): Promise<User> {
	const response = await apiFetch(`/api/auth/bootstrap/${token}`, {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function me(): Promise<User> {
	const response = await apiFetch('/api/auth/me');
	return parseJson(response);
}

export async function updateMe(data: SelfUpdateRequest): Promise<User> {
	const response = await apiFetch('/api/auth/me', {
		method: 'PUT',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function deleteMe(data: SelfDeleteRequest): Promise<void> {
	const response = await apiFetch('/api/auth/me', {
		method: 'DELETE',
		body: JSON.stringify(data)
	});
	await parseJson(response);
}

export async function logout(): Promise<void> {
	const response = await apiFetch('/api/auth/logout', { method: 'POST' }, { minDurationMs: false });
	await parseJson(response);
}
