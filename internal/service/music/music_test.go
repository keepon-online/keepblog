package music

import (
	"testing"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/testutil"
)

func TestGetEnabledMusicList_SortedBySort(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewMusicService(testutil.NewTestDB(t))

	musics, err := s.GetEnabledMusicList()
	if err != nil {
		t.Fatalf("GetEnabledMusicList 报错: %v", err)
	}
	if len(musics) != 2 {
		t.Fatalf("len = %d, want 2（停用曲排除）", len(musics))
	}
	if musics[0].Name != "启用曲A" || musics[1].Name != "启用曲B" {
		t.Errorf("应按 sort 升序: [%s, %s]", musics[0].Name, musics[1].Name)
	}
}

func TestGetMusicList_All(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewMusicService(testutil.NewTestDB(t))

	musics, err := s.GetMusicList()
	if err != nil {
		t.Fatalf("GetMusicList 报错: %v", err)
	}
	if len(musics) != 3 {
		t.Errorf("len = %d, want 3（后台含停用）", len(musics))
	}
}

func TestMusicCRUD(t *testing.T) {
	testutil.NewTestDB(t)
	s := NewMusicService(testutil.NewTestDB(t))

	created := model.Music{Name: "测试曲", Url: "https://cdn.example.com/t.mp3", Sort: 9, State: 1}
	if err := s.Save(created); err != nil {
		t.Fatalf("Save 报错: %v", err)
	}
	list, _ := s.GetMusicList()
	var id uint32
	for _, m := range list {
		if m.Name == "测试曲" {
			id = m.Id
		}
	}
	if id == 0 {
		t.Fatal("新建音乐未找到")
	}

	if err := s.UpdateState(id, 0); err != nil {
		t.Fatalf("UpdateState 报错: %v", err)
	}
	got, _ := s.GetMusic(id)
	if got.State != 0 {
		t.Errorf("State = %d, want 0", got.State)
	}

	if err := s.Delete(int(id)); err != nil {
		t.Fatalf("Delete 报错: %v", err)
	}
	if _, err := s.GetMusic(id); err == nil {
		t.Error("删除后应查不到")
	}
}
