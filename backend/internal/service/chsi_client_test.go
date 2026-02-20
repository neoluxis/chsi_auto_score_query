package service

import (
	"testing"

	"chsi-auto-score-query/pkg/config"
)

func TestParseScoreWithJson(t *testing.T) {
	cfg := &config.Config{
		ChsiUsername: "test",
		ChsiPassword: "test",
	}
	client := NewChsiClient(cfg)

	tests := []struct {
		name        string
		htmlContent string
		expectScore string
		expectError bool
	}{
		{
			name: "Score found",
			htmlContent: `
				<html>
					<script>
						var cj = {"总分": "245", "psyz": "拟录取", "xm": "张三", "zsdwsm": "恭喜被录取"};
						var msg = "success";
					</script>
				</html>
			`,
			expectScore: "245",
			expectError: false,
		},
		{
			name: "Score not published (cj is null)",
			htmlContent: `
				<html>
					<script>
						var cj = null;
						var msg = "成绩尚未发布，请稍后查询";
					</script>
				</html>
			`,
			expectScore: "",
			expectError: false,
		},
		{
			name: "Information mismatch",
			htmlContent: `
				<html>
					<script>
						var cj = null;
						var msg = "信息不匹配，请检查报考信息";
					</script>
				</html>
			`,
			expectScore: "",
			expectError: false,
		},
		{
			name: "Admission status",
			htmlContent: `
				<html>
					<script>
						var cj = {"lqzt": "已录取", "xm": "李四", "ksbh": "103586210002651"};
					</script>
				</html>
			`,
			expectScore: "已录取",
			expectError: false,
		},
		{
			name: "Preliminary score",
			htmlContent: `
				<html>
					<script>
						var cj = {"初试成绩": "385", "xm": "王五", "ksbh": "103586210002652"};
					</script>
				</html>
			`,
			expectScore: "385",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := client.ParseScore(tt.htmlContent)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseScore() error = %v, expectError %v", err, tt.expectError)
			}
			if score == "" && tt.expectScore == "" {
				// Both empty, that's ok
				return
			}
			if score == "" && tt.expectScore != "" {
				t.Errorf("ParseScore() = %q, expected non-empty score", score)
			}
		})
	}
}
