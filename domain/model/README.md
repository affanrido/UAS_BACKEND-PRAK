# Models Directory

Direktori ini berisi semua model data untuk sistem.

## PostgreSQL Models

1. **Users.go** - Model user sistem dengan autentikasi
2. **Roles.go** - Model role/peran user
3. **Permission.go** - Model permission untuk RBAC
4. **Role_Permission.go** - Junction table untuk role-permission
5. **Student.go** - Model data mahasiswa
6. **Lecturers.go** - Model data dosen
7. **achievement_references.go** - Model referensi prestasi (bridge ke MongoDB)

## MongoDB Models

8. **achievements.go** - Model prestasi dinamis mahasiswa

## Struktur Database Hybrid

Project ini menggunakan **hybrid database strategy**:

- **PostgreSQL**: Data relasional dan terstruktur
- **MongoDB**: Data dinamis dengan flexible schema

### Mengapa Hybrid?

**PostgreSQL** untuk:
- Relasi antar entitas (users, roles, students, lecturers)
- Integritas referensial
- ACID transactions
- Data terstruktur

**MongoDB** untuk:
- Data prestasi dengan field yang berbeda-beda per tipe
- Flexible schema (competition, publication, organization, certification)
- Nested objects dan arrays
- Performa query dokumen kompleks

### Bridge Pattern

`achievement_references` di PostgreSQL berfungsi sebagai bridge:
- Menyimpan `mongo_achievement_id` (ObjectId dari MongoDB)
- Menyimpan status verifikasi
- Menjaga relasi dengan student
- Tracking workflow (draft → submitted → verified/rejected)

## Usage

```go
// PostgreSQL Model
import "UAS_BACKEND/domain/model"

user := model.Users{
    Username: "john",
    Email: "john@example.com",
    // ...
}

// MongoDB Model
achievement := model.Achievement{
    StudentID: studentUUID,
    AchievementType: "competition",
    Title: "Hackathon Winner",
    Details: model.AchievementDetails{
        CompetitionName: ptr("National Hackathon 2024"),
        CompetitionLevel: ptr("national"),
        Rank: ptr(1.0),
    },
}
```

## Notes

- Semua model dibuat **tanpa GORM** (menggunakan database/sql untuk PostgreSQL)
- Menggunakan `db` tag untuk mapping kolom database
- Menggunakan `bson` tag untuk MongoDB
- Pointer digunakan untuk optional fields di MongoDB
