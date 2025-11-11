export const TagImages: Record<string, string> = {
    Love: "/images/Love.jpg",
    Study: "/images/Study.jpg",
    Work: "/images/Work.jpg",
    Health: "/images/Health.jpg",
    Finance: "/images/Finance.jpg",
    Personal_Growth: "/images/Personal_Growth.jpg",
  };

export const courseTypes: Record<string, string> = {
  Love: "love",
  Study: "study",
  Work: "work",
  Health: "health",
  Finance: "finance",
  "Personal Growth": "personal_growth",
};

export interface MessageProps {
	ID:        string;
	RoomID:    string;
	SenderID:  string;
	Content:   string;
	Type:      "text" | "notification";
	Status:    "sent" | "delivered" | "read" | "failed";
	CreatedAt: string;
}