// Test translations with deep nesting - Simplified Chinese
const messages = {
  title: 'i18n 多层嵌套测试',
  level1: {
    title: '第一层',
    level2: {
      title: '第二层',
      message: '这是来自第二层的信息',
      level3: {
        title: '第三层',
        message: '这是来自第三层的信息',
        level4: {
          title: '第四层',
          message: '这是来自第四层的信息',
          deepValue: '深度值: {value}',
        }
      }
    }
  }
};

export default messages;

export type TestMessages = typeof messages;
