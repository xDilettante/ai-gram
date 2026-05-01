package telegram

import (
	"encoding/json"
	"testing"
)

func TestMessageDecodesChecklist(t *testing.T) {
	var message Message
	if err := json.Unmarshal([]byte(`{
		"message_id":1,
		"chat":{"id":123,"type":"private"},
		"date":100,
		"reply_to_checklist_task_id":7,
		"checklist":{
			"title":"Release",
			"title_entities":[{"type":"bold","offset":0,"length":7}],
			"tasks":[{"id":1,"text":"Build","text_entities":[{"type":"italic","offset":0,"length":5}],"completed_by_user":{"id":10,"is_bot":false,"first_name":"Ada"},"completion_date":101}],
			"others_can_add_tasks":true,
			"others_can_mark_tasks_as_done":true
		}
	}`), &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.Checklist == nil || message.Checklist.Title != "Release" || len(message.Checklist.Tasks) != 1 {
		t.Fatalf("unexpected checklist: %+v", message.Checklist)
	}
	if message.ReplyToChecklistTaskID != 7 {
		t.Fatalf("unexpected reply_to_checklist_task_id: %d", message.ReplyToChecklistTaskID)
	}
	task := message.Checklist.Tasks[0]
	if task.ID != 1 || task.CompletedByUser == nil || task.CompletedByUser.ID != 10 || task.CompletionDate != 101 {
		t.Fatalf("unexpected task: %+v", task)
	}
}

func TestMessageDecodesChecklistServiceMessages(t *testing.T) {
	var done Message
	if err := json.Unmarshal([]byte(`{
		"message_id":2,
		"chat":{"id":123,"type":"private"},
		"date":100,
		"checklist_tasks_done":{
			"checklist_message":{"message_id":1,"chat":{"id":123,"type":"private"},"date":99},
			"marked_as_done_task_ids":[1,2],
			"marked_as_not_done_task_ids":[3]
		}
	}`), &done); err != nil {
		t.Fatalf("decode done message: %v", err)
	}
	if done.ChecklistTasksDone == nil || len(done.ChecklistTasksDone.MarkedAsDoneTaskIDs) != 2 || len(done.ChecklistTasksDone.MarkedAsNotDoneTaskIDs) != 1 {
		t.Fatalf("unexpected checklist_tasks_done: %+v", done.ChecklistTasksDone)
	}

	var added Message
	if err := json.Unmarshal([]byte(`{
		"message_id":3,
		"chat":{"id":123,"type":"private"},
		"date":100,
		"checklist_tasks_added":{
			"checklist_message":{"message_id":1,"chat":{"id":123,"type":"private"},"date":99},
			"tasks":[{"id":4,"text":"Ship"}]
		}
	}`), &added); err != nil {
		t.Fatalf("decode added message: %v", err)
	}
	if added.ChecklistTasksAdded == nil || len(added.ChecklistTasksAdded.Tasks) != 1 || added.ChecklistTasksAdded.Tasks[0].ID != 4 {
		t.Fatalf("unexpected checklist_tasks_added: %+v", added.ChecklistTasksAdded)
	}
}
