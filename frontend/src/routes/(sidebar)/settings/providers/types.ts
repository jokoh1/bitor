export type ProviderType = 'aws' | 'digitalocean' | 's3' | 'email' | 'slack' | 'teams' | 'discord' | 'telegram' | 'jira';

export interface DigitalOceanSettings {
	region?: string;
	do_project?: string;
	size?: string;
	tags: string[];
	dns_domain?: string;
}

export interface Provider {
	id: string;
	name: string;
	provider_type: ProviderType;
	enabled: boolean;
	uses?: string[];
	settings: DigitalOceanSettings | Record<string, unknown>;
	created?: string;
	updated?: string;
}

export interface ApiKey {
	id: string;
	key: string;
	key_type: string;
	provider: string;
} 