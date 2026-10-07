# Messages

Conversations between the agents of two open stories, one file per conversation, `MS-nnnn-<slug>.md` ([ADR-0120](../../design/adrs/0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md)). Send and answer them with `flai message`, from one story in progress or in review to another; a conversation closes when either story is accepted, cancelled, or archived, and a closed one takes no reply. They are never shown among the operator's threads; `flai check` validates them.
