package entity

type User struct {
	ID               string
	Email            string
	Username         string
	FirstName        string
	Weight           float32
	Height           int
	Age              int
	Sex              string
	PhysicalActivity string
	DayCalories      float32
	Password         string
	BMI              BMIType
}

type BMIType struct {
	Value   float32
	Comment string
}

func (b *BMIType) Calculate(weight float32, height int) {
	b.Value = weight / float32(height*height) * 10000
	switch {
	case b.Value < 16:
		b.Comment = "Выраженный дефицит массы тела"
	case b.Value >= 16 && b.Value < 18.5:
		b.Comment = "Дефицит массы тела"
	case b.Value >= 18.5 && b.Value < 25:
		b.Comment = "Норма"
	case b.Value >= 25 && b.Value < 30:
		b.Comment = "Избыточная масса тела"
	case b.Value >= 30 && b.Value < 35:
		b.Comment = "Ожирение первой степени"
	case b.Value >= 35 && b.Value < 40:
		b.Comment = "Ожирение второй степени"
	case b.Value >= 40:
		b.Comment = "Ожирение третьей степени"
	}
}
