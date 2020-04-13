package valid

import (
	"fmt"
	"time"

	"gopkg.in/go-playground/validator.v10"
)

// https://darjun.github.io/2020/04/04/godailylib/validator/

// Example 下面如未特殊说明，则是根据上面各个类型对应的值与参数值比较
// len：等于参数值，例如len=10；
// max：小于等于参数值，例如max=10；
// min：大于等于参数值，例如min=10；
// eq：等于参数值，注意与len不同。对于字符串，eq约束字符串本身的值，而len约束字符串长度。例如eq=10；
// ne：不等于参数值，例如ne=10；
// gt：大于参数值，例如gt=10；
// gte：大于等于参数值，例如gte=10；
// lt：小于参数值，例如lt=10；
// lte：小于等于参数值，例如lte=10；
// oneof：只能是列举出的值其中一个，这些值必须是数值或字符串，以空格分隔，如果字符串中有空格，将字符串用单引号包围，例如oneof=red green。
// ****** validator中关于字符串的约束有很多，这里介绍几个：
// contains=：包含参数子串，例如contains=email；
// containsany：包含参数中任意的 UNICODE 字符，例如containsany=abcd；
// containsrune：包含参数表示的 rune 字符，例如containsrune=☻；
// excludes：不包含参数子串，例如excludes=email；
// excludesall：不包含参数中任意的 UNICODE 字符，例如excludesall=abcd；
// excludesrune：不包含参数表示的 rune 字符，excludesrune=☻；
// startswith：以参数子串为前缀，例如startswith=hello；
// endswith：以参数子串为后缀，例如endswith=bye。
// ***** 唯一性
// 使用unqiue来指定唯一性约束，对不同类型的处理如下：
// 对于数组和切片，unique约束没有重复的元素；
// 对于map，unique约束没有重复的值；
// 对于元素类型为结构体的切片，unique约束结构体对象的某个字段不重复，通过unqiue=field指定这个字段名。
// ***** 邮件
// 通过email限制字段必须是邮件格式：
// ****** 跨字段约束
// validator允许定义跨字段的约束，即该字段与其他字段之间的关系。
// 这种约束实际上分为两种，一种是参数字段就是同一个结构中的平级字段，另一种是参数字段为结构中其他字段的字段。
// 约束语法很简单，要想使用上面的约束语义，只需要稍微修改一下。
// 例如相等约束（eq），如果是约束同一个结构中的字段，则在后面添加一个field，使用eqfield定义字段间的相等约束。
// 如果是更深层次的字段，在field之前还需要加上cs（可以理解为cross-struct），eq就变为eqcsfield。
// 它们的参数值都是需要比较的字段名，内层的还需要加上字段的类型：
// eqfield=ConfirmPassword
// eqcsfield=InnerStructField.Field
// ***** 特殊
// 有一些比较特殊的约束：
// -：跳过该字段，不检验；
// |：使用多个约束，只需要满足其中一个，例如rgb|rgba；
// required：字段必须设置，不能为默认值；
// omitempty：如果字段未设置，则忽略它。
type Example struct {
	Name      string    `validate:"min=6,max=10"`
	Age       int       `validate:"min=1,max=100"`
	Sex       string    `validate:"oneof=male female"`
	RegTime   time.Time `validate:"lte"`
	Password  string    `validate:"min=10"`
	Password2 string    `validate:"eqfield=Password"`
}

func valid() {
	validate := validator.New()

	u1 := Example{Name: "lidajun", Age: 18}
	err := validate.Struct(u1)
	fmt.Println(err)

	u2 := Example{Name: "dj", Age: 101}
	err = validate.Struct(u2)
	fmt.Println(err)
}
