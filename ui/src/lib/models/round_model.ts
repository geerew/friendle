import { boolean, object, picklist, string, type InferOutput } from 'valibot';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundStatusSchema represents the lifecycle state of a daily round
export const RoundStatusSchema = picklist(['awaiting_word', 'active', 'completed']);

export type RoundStatus = InferOutput<typeof RoundStatusSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundTodaySchema represents today's round state for a group
export const RoundTodaySchema = object({
	roundDate: string(),
	status: RoundStatusSchema,
	isPicker: boolean(),
	roundEndsAt: string()
});

export type RoundTodayModel = InferOutput<typeof RoundTodaySchema>;
