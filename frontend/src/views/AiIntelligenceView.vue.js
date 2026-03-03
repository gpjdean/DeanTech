"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
var vue_1 = require("vue");
var element_plus_1 = require("element-plus");
var icons_vue_1 = require("@element-plus/icons-vue");
var api_1 = require("@/api/api");
// 表单引用
var doubaoForm = (0, vue_1.ref)();
var deepseekForm = (0, vue_1.ref)();
var qianwenForm = (0, vue_1.ref)();
var yuanbaoForm = (0, vue_1.ref)();
var chatMessages = (0, vue_1.ref)(null);
// 选项卡激活状态
var activeTab = (0, vue_1.ref)('doubao');
// 发送状态
var sending = (0, vue_1.ref)(false);
// 选择的AI模型
var selectedModel = (0, vue_1.ref)('doubao');
// 聊天消息
var messages = (0, vue_1.ref)([
    { role: 'assistant', content: '您好，我是AI智能助手，请问有什么可以帮助您的？', timestamp: new Date().toLocaleString() }
]);
// 输入消息
var inputMessage = (0, vue_1.ref)('');
// AI配置数据
var aiConfig = (0, vue_1.reactive)({
    doubao: {
        apiKey: '',
        model: 'doubao-pro-4',
        temperature: 0.7
    },
    deepseek: {
        apiKey: '',
        model: 'deepseek-r1',
        temperature: 0.7
    },
    qianwen: {
        apiKey: '',
        model: 'qianwen-4',
        temperature: 0.7
    },
    yuanbao: {
        apiKey: '',
        model: 'yuanbao-4',
        temperature: 0.7
    }
});
// 组件挂载时加载配置
(0, vue_1.onMounted)(function () {
    loadConfig();
});
// 加载配置
var loadConfig = function () { return __awaiter(void 0, void 0, void 0, function () {
    var response, settings, aiSetting, config, error_1;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                _a.trys.push([0, 2, , 3]);
                return [4 /*yield*/, api_1.systemSettingAPI.list()];
            case 1:
                response = _a.sent();
                settings = response.data;
                aiSetting = settings.find(function (setting) { return setting.key === 'aiConfig'; });
                if (aiSetting) {
                    config = JSON.parse(aiSetting.value);
                    Object.assign(aiConfig, config);
                }
                return [3 /*break*/, 3];
            case 2:
                error_1 = _a.sent();
                console.error('加载AI配置失败:', error_1);
                element_plus_1.ElMessage.error('加载AI配置失败');
                return [3 /*break*/, 3];
            case 3: return [2 /*return*/];
        }
    });
}); };
// 保存配置
var saveConfig = function (model) { return __awaiter(void 0, void 0, void 0, function () {
    var configData, error_2, error_3;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                _a.trys.push([0, 8, , 9]);
                configData = {
                    key: 'aiConfig',
                    value: JSON.stringify(aiConfig),
                    description: 'AI智能助手配置'
                };
                _a.label = 1;
            case 1:
                _a.trys.push([1, 3, , 7]);
                // 尝试更新现有配置
                return [4 /*yield*/, api_1.systemSettingAPI.update('aiConfig', configData)];
            case 2:
                // 尝试更新现有配置
                _a.sent();
                return [3 /*break*/, 7];
            case 3:
                error_2 = _a.sent();
                if (!(error_2.response && error_2.response.status === 404)) return [3 /*break*/, 5];
                return [4 /*yield*/, api_1.systemSettingAPI.create(configData)];
            case 4:
                _a.sent();
                return [3 /*break*/, 6];
            case 5: throw error_2;
            case 6: return [3 /*break*/, 7];
            case 7:
                element_plus_1.ElMessage.success("".concat(model, "\u914D\u7F6E\u4FDD\u5B58\u6210\u529F"));
                return [3 /*break*/, 9];
            case 8:
                error_3 = _a.sent();
                console.error('保存AI配置失败:', error_3);
                element_plus_1.ElMessage.error('保存AI配置失败');
                return [3 /*break*/, 9];
            case 9: return [2 /*return*/];
        }
    });
}); };
// 测试连接
var testConnection = function (model) { return __awaiter(void 0, void 0, void 0, function () {
    var response, error_4;
    var _a, _b;
    return __generator(this, function (_c) {
        switch (_c.label) {
            case 0:
                _c.trys.push([0, 2, , 3]);
                element_plus_1.ElMessage.info("\u6B63\u5728\u6D4B\u8BD5".concat(model, "\u8FDE\u63A5..."));
                return [4 /*yield*/, api_1.aiAPI.testConnection({
                        model: model,
                        apiKey: aiConfig[model].apiKey
                    })];
            case 1:
                response = _c.sent();
                (0, element_plus_1.ElNotification)({
                    title: '连接测试',
                    message: response.data.message,
                    type: 'success'
                });
                return [3 /*break*/, 3];
            case 2:
                error_4 = _c.sent();
                console.error('测试连接失败:', error_4);
                element_plus_1.ElMessage.error(((_b = (_a = error_4.response) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.error) || '测试连接失败');
                return [3 /*break*/, 3];
            case 3: return [2 /*return*/];
        }
    });
}); };
// 发送消息
var sendMessage = function () { return __awaiter(void 0, void 0, void 0, function () {
    var userMessage, chatHistory, response, aiMessage, error_5;
    var _a, _b;
    return __generator(this, function (_c) {
        switch (_c.label) {
            case 0:
                if (!inputMessage.value.trim()) {
                    element_plus_1.ElMessage.warning('请输入消息内容');
                    return [2 /*return*/];
                }
                userMessage = {
                    role: 'user',
                    content: inputMessage.value,
                    timestamp: new Date().toLocaleString()
                };
                messages.value.push(userMessage);
                inputMessage.value = '';
                // 滚动到底部
                return [4 /*yield*/, (0, vue_1.nextTick)()];
            case 1:
                // 滚动到底部
                _c.sent();
                scrollToBottom();
                // 发送消息到后端获取AI回复
                sending.value = true;
                _c.label = 2;
            case 2:
                _c.trys.push([2, 4, , 5]);
                chatHistory = messages.value.map(function (msg) { return ({
                    role: msg.role,
                    content: msg.content
                }); });
                return [4 /*yield*/, api_1.aiAPI.chat({
                        model: selectedModel.value,
                        message: userMessage.content,
                        history: chatHistory
                    })
                    // 添加AI回复
                ];
            case 3:
                response = _c.sent();
                aiMessage = {
                    role: 'assistant',
                    content: response.data.response,
                    timestamp: new Date().toLocaleString()
                };
                messages.value.push(aiMessage);
                sending.value = false;
                scrollToBottom();
                return [3 /*break*/, 5];
            case 4:
                error_5 = _c.sent();
                console.error('发送消息失败:', error_5);
                element_plus_1.ElMessage.error(((_b = (_a = error_5.response) === null || _a === void 0 ? void 0 : _a.data) === null || _b === void 0 ? void 0 : _b.error) || '发送消息失败');
                sending.value = false;
                return [3 /*break*/, 5];
            case 5: return [2 /*return*/];
        }
    });
}); };
// 清空聊天
var clearChat = function () {
    messages.value = [
        { role: 'assistant', content: '您好，我是AI智能助手，请问有什么可以帮助您的？', timestamp: new Date().toLocaleString() }
    ];
};
// 滚动到底部
var scrollToBottom = function () {
    if (chatMessages.value) {
        chatMessages.value.scrollTop = chatMessages.value.scrollHeight;
    }
};
debugger; /* PartiallyEnd: #3632/scriptSetup.vue */
var __VLS_ctx = {};
var __VLS_components;
var __VLS_directives;
/** @type {__VLS_StyleScopedClasses['header-content']} */ ;
/** @type {__VLS_StyleScopedClasses['message-item']} */ ;
/** @type {__VLS_StyleScopedClasses['message-item']} */ ;
/** @type {__VLS_StyleScopedClasses['user']} */ ;
/** @type {__VLS_StyleScopedClasses['message-text']} */ ;
/** @type {__VLS_StyleScopedClasses['message-item']} */ ;
/** @type {__VLS_StyleScopedClasses['message-text']} */ ;
/** @type {__VLS_StyleScopedClasses['message-item']} */ ;
/** @type {__VLS_StyleScopedClasses['assistant']} */ ;
/** @type {__VLS_StyleScopedClasses['message-time']} */ ;
// CSS variable injection 
// CSS variable injection end 
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "ai-intelligence-view" }));
var __VLS_0 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
var __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0(__assign({ shadow: "never" }, { class: "page-header" })));
var __VLS_2 = __VLS_1.apply(void 0, __spreadArray([__assign({ shadow: "never" }, { class: "page-header" })], __VLS_functionalComponentArgsRest(__VLS_1), false));
__VLS_3.slots.default;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "header-content" }));
__VLS_asFunctionalElement(__VLS_intrinsicElements.h2, __VLS_intrinsicElements.h2)({});
__VLS_asFunctionalElement(__VLS_intrinsicElements.p, __VLS_intrinsicElements.p)({});
var __VLS_3;
var __VLS_4 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
var __VLS_5 = __VLS_asFunctionalComponent(__VLS_4, new __VLS_4(__assign({ shadow: "hover" }, { class: "config-card" })));
var __VLS_6 = __VLS_5.apply(void 0, __spreadArray([__assign({ shadow: "hover" }, { class: "config-card" })], __VLS_functionalComponentArgsRest(__VLS_5), false));
__VLS_7.slots.default;
{
    var __VLS_thisSlot = __VLS_7.slots.header;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "card-header" }));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
}
var __VLS_8 = {}.ElTabs;
/** @type {[typeof __VLS_components.ElTabs, typeof __VLS_components.elTabs, typeof __VLS_components.ElTabs, typeof __VLS_components.elTabs, ]} */ ;
// @ts-ignore
var __VLS_9 = __VLS_asFunctionalComponent(__VLS_8, new __VLS_8(__assign({ modelValue: (__VLS_ctx.activeTab) }, { class: "ai-tabs" })));
var __VLS_10 = __VLS_9.apply(void 0, __spreadArray([__assign({ modelValue: (__VLS_ctx.activeTab) }, { class: "ai-tabs" })], __VLS_functionalComponentArgsRest(__VLS_9), false));
__VLS_11.slots.default;
var __VLS_12 = {}.ElTabPane;
/** @type {[typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, ]} */ ;
// @ts-ignore
var __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
    label: "豆包",
    name: "doubao",
}));
var __VLS_14 = __VLS_13.apply(void 0, __spreadArray([{
        label: "豆包",
        name: "doubao",
    }], __VLS_functionalComponentArgsRest(__VLS_13), false));
__VLS_15.slots.default;
var __VLS_16 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
var __VLS_17 = __VLS_asFunctionalComponent(__VLS_16, new __VLS_16(__assign({ ref: "doubaoForm", model: (__VLS_ctx.aiConfig.doubao), labelWidth: "120px" }, { class: "ai-form" })));
var __VLS_18 = __VLS_17.apply(void 0, __spreadArray([__assign({ ref: "doubaoForm", model: (__VLS_ctx.aiConfig.doubao), labelWidth: "120px" }, { class: "ai-form" })], __VLS_functionalComponentArgsRest(__VLS_17), false));
/** @type {typeof __VLS_ctx.doubaoForm} */ ;
var __VLS_20 = {};
__VLS_19.slots.default;
var __VLS_22 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_23 = __VLS_asFunctionalComponent(__VLS_22, new __VLS_22({
    label: "API Key",
    prop: "apiKey",
}));
var __VLS_24 = __VLS_23.apply(void 0, __spreadArray([{
        label: "API Key",
        prop: "apiKey",
    }], __VLS_functionalComponentArgsRest(__VLS_23), false));
__VLS_25.slots.default;
var __VLS_26 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
var __VLS_27 = __VLS_asFunctionalComponent(__VLS_26, new __VLS_26({
    modelValue: (__VLS_ctx.aiConfig.doubao.apiKey),
    type: "password",
    placeholder: "请输入豆包API Key",
    showPassword: true,
}));
var __VLS_28 = __VLS_27.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.doubao.apiKey),
        type: "password",
        placeholder: "请输入豆包API Key",
        showPassword: true,
    }], __VLS_functionalComponentArgsRest(__VLS_27), false));
var __VLS_25;
var __VLS_30 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_31 = __VLS_asFunctionalComponent(__VLS_30, new __VLS_30({
    label: "模型名称",
    prop: "model",
}));
var __VLS_32 = __VLS_31.apply(void 0, __spreadArray([{
        label: "模型名称",
        prop: "model",
    }], __VLS_functionalComponentArgsRest(__VLS_31), false));
__VLS_33.slots.default;
var __VLS_34 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
var __VLS_35 = __VLS_asFunctionalComponent(__VLS_34, new __VLS_34({
    modelValue: (__VLS_ctx.aiConfig.doubao.model),
    placeholder: "请选择豆包模型",
}));
var __VLS_36 = __VLS_35.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.doubao.model),
        placeholder: "请选择豆包模型",
    }], __VLS_functionalComponentArgsRest(__VLS_35), false));
__VLS_37.slots.default;
var __VLS_38 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_39 = __VLS_asFunctionalComponent(__VLS_38, new __VLS_38({
    label: "豆包-4",
    value: "doubao-pro-4",
}));
var __VLS_40 = __VLS_39.apply(void 0, __spreadArray([{
        label: "豆包-4",
        value: "doubao-pro-4",
    }], __VLS_functionalComponentArgsRest(__VLS_39), false));
var __VLS_42 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_43 = __VLS_asFunctionalComponent(__VLS_42, new __VLS_42({
    label: "豆包-3",
    value: "doubao-pro-3",
}));
var __VLS_44 = __VLS_43.apply(void 0, __spreadArray([{
        label: "豆包-3",
        value: "doubao-pro-3",
    }], __VLS_functionalComponentArgsRest(__VLS_43), false));
var __VLS_46 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_47 = __VLS_asFunctionalComponent(__VLS_46, new __VLS_46({
    label: "豆包-2",
    value: "doubao-pro-2",
}));
var __VLS_48 = __VLS_47.apply(void 0, __spreadArray([{
        label: "豆包-2",
        value: "doubao-pro-2",
    }], __VLS_functionalComponentArgsRest(__VLS_47), false));
var __VLS_37;
var __VLS_33;
var __VLS_50 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_51 = __VLS_asFunctionalComponent(__VLS_50, new __VLS_50({
    label: "温度",
    prop: "temperature",
}));
var __VLS_52 = __VLS_51.apply(void 0, __spreadArray([{
        label: "温度",
        prop: "temperature",
    }], __VLS_functionalComponentArgsRest(__VLS_51), false));
__VLS_53.slots.default;
var __VLS_54 = {}.ElSlider;
/** @type {[typeof __VLS_components.ElSlider, typeof __VLS_components.elSlider, ]} */ ;
// @ts-ignore
var __VLS_55 = __VLS_asFunctionalComponent(__VLS_54, new __VLS_54({
    modelValue: (__VLS_ctx.aiConfig.doubao.temperature),
    min: (0),
    max: (1),
    step: (0.1),
}));
var __VLS_56 = __VLS_55.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.doubao.temperature),
        min: (0),
        max: (1),
        step: (0.1),
    }], __VLS_functionalComponentArgsRest(__VLS_55), false));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)(__assign({ class: "slider-value" }));
(__VLS_ctx.aiConfig.doubao.temperature);
var __VLS_53;
var __VLS_58 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_59 = __VLS_asFunctionalComponent(__VLS_58, new __VLS_58({}));
var __VLS_60 = __VLS_59.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_59), false));
__VLS_61.slots.default;
var __VLS_62 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_63 = __VLS_asFunctionalComponent(__VLS_62, new __VLS_62(__assign({ 'onClick': {} }, { type: "primary" })));
var __VLS_64 = __VLS_63.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { type: "primary" })], __VLS_functionalComponentArgsRest(__VLS_63), false));
var __VLS_66;
var __VLS_67;
var __VLS_68;
var __VLS_69 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.saveConfig('doubao');
    }
};
__VLS_65.slots.default;
var __VLS_65;
var __VLS_70 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_71 = __VLS_asFunctionalComponent(__VLS_70, new __VLS_70(__assign({ 'onClick': {} })));
var __VLS_72 = __VLS_71.apply(void 0, __spreadArray([__assign({ 'onClick': {} })], __VLS_functionalComponentArgsRest(__VLS_71), false));
var __VLS_74;
var __VLS_75;
var __VLS_76;
var __VLS_77 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.testConnection('doubao');
    }
};
__VLS_73.slots.default;
var __VLS_73;
var __VLS_61;
var __VLS_19;
var __VLS_15;
var __VLS_78 = {}.ElTabPane;
/** @type {[typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, ]} */ ;
// @ts-ignore
var __VLS_79 = __VLS_asFunctionalComponent(__VLS_78, new __VLS_78({
    label: "DeepSeek",
    name: "deepseek",
}));
var __VLS_80 = __VLS_79.apply(void 0, __spreadArray([{
        label: "DeepSeek",
        name: "deepseek",
    }], __VLS_functionalComponentArgsRest(__VLS_79), false));
__VLS_81.slots.default;
var __VLS_82 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
var __VLS_83 = __VLS_asFunctionalComponent(__VLS_82, new __VLS_82(__assign({ ref: "deepseekForm", model: (__VLS_ctx.aiConfig.deepseek), labelWidth: "120px" }, { class: "ai-form" })));
var __VLS_84 = __VLS_83.apply(void 0, __spreadArray([__assign({ ref: "deepseekForm", model: (__VLS_ctx.aiConfig.deepseek), labelWidth: "120px" }, { class: "ai-form" })], __VLS_functionalComponentArgsRest(__VLS_83), false));
/** @type {typeof __VLS_ctx.deepseekForm} */ ;
var __VLS_86 = {};
__VLS_85.slots.default;
var __VLS_88 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_89 = __VLS_asFunctionalComponent(__VLS_88, new __VLS_88({
    label: "API Key",
    prop: "apiKey",
}));
var __VLS_90 = __VLS_89.apply(void 0, __spreadArray([{
        label: "API Key",
        prop: "apiKey",
    }], __VLS_functionalComponentArgsRest(__VLS_89), false));
__VLS_91.slots.default;
var __VLS_92 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
var __VLS_93 = __VLS_asFunctionalComponent(__VLS_92, new __VLS_92({
    modelValue: (__VLS_ctx.aiConfig.deepseek.apiKey),
    type: "password",
    placeholder: "请输入DeepSeek API Key",
    showPassword: true,
}));
var __VLS_94 = __VLS_93.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.deepseek.apiKey),
        type: "password",
        placeholder: "请输入DeepSeek API Key",
        showPassword: true,
    }], __VLS_functionalComponentArgsRest(__VLS_93), false));
var __VLS_91;
var __VLS_96 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_97 = __VLS_asFunctionalComponent(__VLS_96, new __VLS_96({
    label: "模型名称",
    prop: "model",
}));
var __VLS_98 = __VLS_97.apply(void 0, __spreadArray([{
        label: "模型名称",
        prop: "model",
    }], __VLS_functionalComponentArgsRest(__VLS_97), false));
__VLS_99.slots.default;
var __VLS_100 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
var __VLS_101 = __VLS_asFunctionalComponent(__VLS_100, new __VLS_100({
    modelValue: (__VLS_ctx.aiConfig.deepseek.model),
    placeholder: "请选择DeepSeek模型",
}));
var __VLS_102 = __VLS_101.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.deepseek.model),
        placeholder: "请选择DeepSeek模型",
    }], __VLS_functionalComponentArgsRest(__VLS_101), false));
__VLS_103.slots.default;
var __VLS_104 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_105 = __VLS_asFunctionalComponent(__VLS_104, new __VLS_104({
    label: "DeepSeek-R1",
    value: "deepseek-r1",
}));
var __VLS_106 = __VLS_105.apply(void 0, __spreadArray([{
        label: "DeepSeek-R1",
        value: "deepseek-r1",
    }], __VLS_functionalComponentArgsRest(__VLS_105), false));
var __VLS_108 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_109 = __VLS_asFunctionalComponent(__VLS_108, new __VLS_108({
    label: "DeepSeek-V2",
    value: "deepseek-v2",
}));
var __VLS_110 = __VLS_109.apply(void 0, __spreadArray([{
        label: "DeepSeek-V2",
        value: "deepseek-v2",
    }], __VLS_functionalComponentArgsRest(__VLS_109), false));
var __VLS_112 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_113 = __VLS_asFunctionalComponent(__VLS_112, new __VLS_112({
    label: "DeepSeek-Coder",
    value: "deepseek-coder",
}));
var __VLS_114 = __VLS_113.apply(void 0, __spreadArray([{
        label: "DeepSeek-Coder",
        value: "deepseek-coder",
    }], __VLS_functionalComponentArgsRest(__VLS_113), false));
var __VLS_103;
var __VLS_99;
var __VLS_116 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_117 = __VLS_asFunctionalComponent(__VLS_116, new __VLS_116({
    label: "温度",
    prop: "temperature",
}));
var __VLS_118 = __VLS_117.apply(void 0, __spreadArray([{
        label: "温度",
        prop: "temperature",
    }], __VLS_functionalComponentArgsRest(__VLS_117), false));
__VLS_119.slots.default;
var __VLS_120 = {}.ElSlider;
/** @type {[typeof __VLS_components.ElSlider, typeof __VLS_components.elSlider, ]} */ ;
// @ts-ignore
var __VLS_121 = __VLS_asFunctionalComponent(__VLS_120, new __VLS_120({
    modelValue: (__VLS_ctx.aiConfig.deepseek.temperature),
    min: (0),
    max: (1),
    step: (0.1),
}));
var __VLS_122 = __VLS_121.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.deepseek.temperature),
        min: (0),
        max: (1),
        step: (0.1),
    }], __VLS_functionalComponentArgsRest(__VLS_121), false));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)(__assign({ class: "slider-value" }));
(__VLS_ctx.aiConfig.deepseek.temperature);
var __VLS_119;
var __VLS_124 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_125 = __VLS_asFunctionalComponent(__VLS_124, new __VLS_124({}));
var __VLS_126 = __VLS_125.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_125), false));
__VLS_127.slots.default;
var __VLS_128 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_129 = __VLS_asFunctionalComponent(__VLS_128, new __VLS_128(__assign({ 'onClick': {} }, { type: "primary" })));
var __VLS_130 = __VLS_129.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { type: "primary" })], __VLS_functionalComponentArgsRest(__VLS_129), false));
var __VLS_132;
var __VLS_133;
var __VLS_134;
var __VLS_135 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.saveConfig('deepseek');
    }
};
__VLS_131.slots.default;
var __VLS_131;
var __VLS_136 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_137 = __VLS_asFunctionalComponent(__VLS_136, new __VLS_136(__assign({ 'onClick': {} })));
var __VLS_138 = __VLS_137.apply(void 0, __spreadArray([__assign({ 'onClick': {} })], __VLS_functionalComponentArgsRest(__VLS_137), false));
var __VLS_140;
var __VLS_141;
var __VLS_142;
var __VLS_143 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.testConnection('deepseek');
    }
};
__VLS_139.slots.default;
var __VLS_139;
var __VLS_127;
var __VLS_85;
var __VLS_81;
var __VLS_144 = {}.ElTabPane;
/** @type {[typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, ]} */ ;
// @ts-ignore
var __VLS_145 = __VLS_asFunctionalComponent(__VLS_144, new __VLS_144({
    label: "千问",
    name: "qianwen",
}));
var __VLS_146 = __VLS_145.apply(void 0, __spreadArray([{
        label: "千问",
        name: "qianwen",
    }], __VLS_functionalComponentArgsRest(__VLS_145), false));
__VLS_147.slots.default;
var __VLS_148 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
var __VLS_149 = __VLS_asFunctionalComponent(__VLS_148, new __VLS_148(__assign({ ref: "qianwenForm", model: (__VLS_ctx.aiConfig.qianwen), labelWidth: "120px" }, { class: "ai-form" })));
var __VLS_150 = __VLS_149.apply(void 0, __spreadArray([__assign({ ref: "qianwenForm", model: (__VLS_ctx.aiConfig.qianwen), labelWidth: "120px" }, { class: "ai-form" })], __VLS_functionalComponentArgsRest(__VLS_149), false));
/** @type {typeof __VLS_ctx.qianwenForm} */ ;
var __VLS_152 = {};
__VLS_151.slots.default;
var __VLS_154 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_155 = __VLS_asFunctionalComponent(__VLS_154, new __VLS_154({
    label: "API Key",
    prop: "apiKey",
}));
var __VLS_156 = __VLS_155.apply(void 0, __spreadArray([{
        label: "API Key",
        prop: "apiKey",
    }], __VLS_functionalComponentArgsRest(__VLS_155), false));
__VLS_157.slots.default;
var __VLS_158 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
var __VLS_159 = __VLS_asFunctionalComponent(__VLS_158, new __VLS_158({
    modelValue: (__VLS_ctx.aiConfig.qianwen.apiKey),
    type: "password",
    placeholder: "请输入千问API Key",
    showPassword: true,
}));
var __VLS_160 = __VLS_159.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.qianwen.apiKey),
        type: "password",
        placeholder: "请输入千问API Key",
        showPassword: true,
    }], __VLS_functionalComponentArgsRest(__VLS_159), false));
var __VLS_157;
var __VLS_162 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_163 = __VLS_asFunctionalComponent(__VLS_162, new __VLS_162({
    label: "模型名称",
    prop: "model",
}));
var __VLS_164 = __VLS_163.apply(void 0, __spreadArray([{
        label: "模型名称",
        prop: "model",
    }], __VLS_functionalComponentArgsRest(__VLS_163), false));
__VLS_165.slots.default;
var __VLS_166 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
var __VLS_167 = __VLS_asFunctionalComponent(__VLS_166, new __VLS_166({
    modelValue: (__VLS_ctx.aiConfig.qianwen.model),
    placeholder: "请选择千问模型",
}));
var __VLS_168 = __VLS_167.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.qianwen.model),
        placeholder: "请选择千问模型",
    }], __VLS_functionalComponentArgsRest(__VLS_167), false));
__VLS_169.slots.default;
var __VLS_170 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_171 = __VLS_asFunctionalComponent(__VLS_170, new __VLS_170({
    label: "千问-4",
    value: "qianwen-4",
}));
var __VLS_172 = __VLS_171.apply(void 0, __spreadArray([{
        label: "千问-4",
        value: "qianwen-4",
    }], __VLS_functionalComponentArgsRest(__VLS_171), false));
var __VLS_174 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_175 = __VLS_asFunctionalComponent(__VLS_174, new __VLS_174({
    label: "千问-3",
    value: "qianwen-3",
}));
var __VLS_176 = __VLS_175.apply(void 0, __spreadArray([{
        label: "千问-3",
        value: "qianwen-3",
    }], __VLS_functionalComponentArgsRest(__VLS_175), false));
var __VLS_178 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_179 = __VLS_asFunctionalComponent(__VLS_178, new __VLS_178({
    label: "千问-2",
    value: "qianwen-2",
}));
var __VLS_180 = __VLS_179.apply(void 0, __spreadArray([{
        label: "千问-2",
        value: "qianwen-2",
    }], __VLS_functionalComponentArgsRest(__VLS_179), false));
var __VLS_169;
var __VLS_165;
var __VLS_182 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_183 = __VLS_asFunctionalComponent(__VLS_182, new __VLS_182({
    label: "温度",
    prop: "temperature",
}));
var __VLS_184 = __VLS_183.apply(void 0, __spreadArray([{
        label: "温度",
        prop: "temperature",
    }], __VLS_functionalComponentArgsRest(__VLS_183), false));
__VLS_185.slots.default;
var __VLS_186 = {}.ElSlider;
/** @type {[typeof __VLS_components.ElSlider, typeof __VLS_components.elSlider, ]} */ ;
// @ts-ignore
var __VLS_187 = __VLS_asFunctionalComponent(__VLS_186, new __VLS_186({
    modelValue: (__VLS_ctx.aiConfig.qianwen.temperature),
    min: (0),
    max: (1),
    step: (0.1),
}));
var __VLS_188 = __VLS_187.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.qianwen.temperature),
        min: (0),
        max: (1),
        step: (0.1),
    }], __VLS_functionalComponentArgsRest(__VLS_187), false));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)(__assign({ class: "slider-value" }));
(__VLS_ctx.aiConfig.qianwen.temperature);
var __VLS_185;
var __VLS_190 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_191 = __VLS_asFunctionalComponent(__VLS_190, new __VLS_190({}));
var __VLS_192 = __VLS_191.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_191), false));
__VLS_193.slots.default;
var __VLS_194 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_195 = __VLS_asFunctionalComponent(__VLS_194, new __VLS_194(__assign({ 'onClick': {} }, { type: "primary" })));
var __VLS_196 = __VLS_195.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { type: "primary" })], __VLS_functionalComponentArgsRest(__VLS_195), false));
var __VLS_198;
var __VLS_199;
var __VLS_200;
var __VLS_201 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.saveConfig('qianwen');
    }
};
__VLS_197.slots.default;
var __VLS_197;
var __VLS_202 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_203 = __VLS_asFunctionalComponent(__VLS_202, new __VLS_202(__assign({ 'onClick': {} })));
var __VLS_204 = __VLS_203.apply(void 0, __spreadArray([__assign({ 'onClick': {} })], __VLS_functionalComponentArgsRest(__VLS_203), false));
var __VLS_206;
var __VLS_207;
var __VLS_208;
var __VLS_209 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.testConnection('qianwen');
    }
};
__VLS_205.slots.default;
var __VLS_205;
var __VLS_193;
var __VLS_151;
var __VLS_147;
var __VLS_210 = {}.ElTabPane;
/** @type {[typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, typeof __VLS_components.ElTabPane, typeof __VLS_components.elTabPane, ]} */ ;
// @ts-ignore
var __VLS_211 = __VLS_asFunctionalComponent(__VLS_210, new __VLS_210({
    label: "元宝",
    name: "yuanbao",
}));
var __VLS_212 = __VLS_211.apply(void 0, __spreadArray([{
        label: "元宝",
        name: "yuanbao",
    }], __VLS_functionalComponentArgsRest(__VLS_211), false));
__VLS_213.slots.default;
var __VLS_214 = {}.ElForm;
/** @type {[typeof __VLS_components.ElForm, typeof __VLS_components.elForm, typeof __VLS_components.ElForm, typeof __VLS_components.elForm, ]} */ ;
// @ts-ignore
var __VLS_215 = __VLS_asFunctionalComponent(__VLS_214, new __VLS_214(__assign({ ref: "yuanbaoForm", model: (__VLS_ctx.aiConfig.yuanbao), labelWidth: "120px" }, { class: "ai-form" })));
var __VLS_216 = __VLS_215.apply(void 0, __spreadArray([__assign({ ref: "yuanbaoForm", model: (__VLS_ctx.aiConfig.yuanbao), labelWidth: "120px" }, { class: "ai-form" })], __VLS_functionalComponentArgsRest(__VLS_215), false));
/** @type {typeof __VLS_ctx.yuanbaoForm} */ ;
var __VLS_218 = {};
__VLS_217.slots.default;
var __VLS_220 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_221 = __VLS_asFunctionalComponent(__VLS_220, new __VLS_220({
    label: "API Key",
    prop: "apiKey",
}));
var __VLS_222 = __VLS_221.apply(void 0, __spreadArray([{
        label: "API Key",
        prop: "apiKey",
    }], __VLS_functionalComponentArgsRest(__VLS_221), false));
__VLS_223.slots.default;
var __VLS_224 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
var __VLS_225 = __VLS_asFunctionalComponent(__VLS_224, new __VLS_224({
    modelValue: (__VLS_ctx.aiConfig.yuanbao.apiKey),
    type: "password",
    placeholder: "请输入元宝API Key",
    showPassword: true,
}));
var __VLS_226 = __VLS_225.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.yuanbao.apiKey),
        type: "password",
        placeholder: "请输入元宝API Key",
        showPassword: true,
    }], __VLS_functionalComponentArgsRest(__VLS_225), false));
var __VLS_223;
var __VLS_228 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_229 = __VLS_asFunctionalComponent(__VLS_228, new __VLS_228({
    label: "模型名称",
    prop: "model",
}));
var __VLS_230 = __VLS_229.apply(void 0, __spreadArray([{
        label: "模型名称",
        prop: "model",
    }], __VLS_functionalComponentArgsRest(__VLS_229), false));
__VLS_231.slots.default;
var __VLS_232 = {}.ElSelect;
/** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
// @ts-ignore
var __VLS_233 = __VLS_asFunctionalComponent(__VLS_232, new __VLS_232({
    modelValue: (__VLS_ctx.aiConfig.yuanbao.model),
    placeholder: "请选择元宝模型",
}));
var __VLS_234 = __VLS_233.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.yuanbao.model),
        placeholder: "请选择元宝模型",
    }], __VLS_functionalComponentArgsRest(__VLS_233), false));
__VLS_235.slots.default;
var __VLS_236 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_237 = __VLS_asFunctionalComponent(__VLS_236, new __VLS_236({
    label: "元宝-4",
    value: "yuanbao-4",
}));
var __VLS_238 = __VLS_237.apply(void 0, __spreadArray([{
        label: "元宝-4",
        value: "yuanbao-4",
    }], __VLS_functionalComponentArgsRest(__VLS_237), false));
var __VLS_240 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_241 = __VLS_asFunctionalComponent(__VLS_240, new __VLS_240({
    label: "元宝-3",
    value: "yuanbao-3",
}));
var __VLS_242 = __VLS_241.apply(void 0, __spreadArray([{
        label: "元宝-3",
        value: "yuanbao-3",
    }], __VLS_functionalComponentArgsRest(__VLS_241), false));
var __VLS_244 = {}.ElOption;
/** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
// @ts-ignore
var __VLS_245 = __VLS_asFunctionalComponent(__VLS_244, new __VLS_244({
    label: "元宝-2",
    value: "yuanbao-2",
}));
var __VLS_246 = __VLS_245.apply(void 0, __spreadArray([{
        label: "元宝-2",
        value: "yuanbao-2",
    }], __VLS_functionalComponentArgsRest(__VLS_245), false));
var __VLS_235;
var __VLS_231;
var __VLS_248 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_249 = __VLS_asFunctionalComponent(__VLS_248, new __VLS_248({
    label: "温度",
    prop: "temperature",
}));
var __VLS_250 = __VLS_249.apply(void 0, __spreadArray([{
        label: "温度",
        prop: "temperature",
    }], __VLS_functionalComponentArgsRest(__VLS_249), false));
__VLS_251.slots.default;
var __VLS_252 = {}.ElSlider;
/** @type {[typeof __VLS_components.ElSlider, typeof __VLS_components.elSlider, ]} */ ;
// @ts-ignore
var __VLS_253 = __VLS_asFunctionalComponent(__VLS_252, new __VLS_252({
    modelValue: (__VLS_ctx.aiConfig.yuanbao.temperature),
    min: (0),
    max: (1),
    step: (0.1),
}));
var __VLS_254 = __VLS_253.apply(void 0, __spreadArray([{
        modelValue: (__VLS_ctx.aiConfig.yuanbao.temperature),
        min: (0),
        max: (1),
        step: (0.1),
    }], __VLS_functionalComponentArgsRest(__VLS_253), false));
__VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)(__assign({ class: "slider-value" }));
(__VLS_ctx.aiConfig.yuanbao.temperature);
var __VLS_251;
var __VLS_256 = {}.ElFormItem;
/** @type {[typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, typeof __VLS_components.ElFormItem, typeof __VLS_components.elFormItem, ]} */ ;
// @ts-ignore
var __VLS_257 = __VLS_asFunctionalComponent(__VLS_256, new __VLS_256({}));
var __VLS_258 = __VLS_257.apply(void 0, __spreadArray([{}], __VLS_functionalComponentArgsRest(__VLS_257), false));
__VLS_259.slots.default;
var __VLS_260 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_261 = __VLS_asFunctionalComponent(__VLS_260, new __VLS_260(__assign({ 'onClick': {} }, { type: "primary" })));
var __VLS_262 = __VLS_261.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { type: "primary" })], __VLS_functionalComponentArgsRest(__VLS_261), false));
var __VLS_264;
var __VLS_265;
var __VLS_266;
var __VLS_267 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.saveConfig('yuanbao');
    }
};
__VLS_263.slots.default;
var __VLS_263;
var __VLS_268 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_269 = __VLS_asFunctionalComponent(__VLS_268, new __VLS_268(__assign({ 'onClick': {} })));
var __VLS_270 = __VLS_269.apply(void 0, __spreadArray([__assign({ 'onClick': {} })], __VLS_functionalComponentArgsRest(__VLS_269), false));
var __VLS_272;
var __VLS_273;
var __VLS_274;
var __VLS_275 = {
    onClick: function () {
        var _a = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            _a[_i] = arguments[_i];
        }
        var $event = _a[0];
        __VLS_ctx.testConnection('yuanbao');
    }
};
__VLS_271.slots.default;
var __VLS_271;
var __VLS_259;
var __VLS_217;
var __VLS_213;
var __VLS_11;
var __VLS_7;
var __VLS_276 = {}.ElCard;
/** @type {[typeof __VLS_components.ElCard, typeof __VLS_components.elCard, typeof __VLS_components.ElCard, typeof __VLS_components.elCard, ]} */ ;
// @ts-ignore
var __VLS_277 = __VLS_asFunctionalComponent(__VLS_276, new __VLS_276(__assign({ shadow: "hover" }, { class: "chat-card" })));
var __VLS_278 = __VLS_277.apply(void 0, __spreadArray([__assign({ shadow: "hover" }, { class: "chat-card" })], __VLS_functionalComponentArgsRest(__VLS_277), false));
__VLS_279.slots.default;
{
    var __VLS_thisSlot = __VLS_279.slots.header;
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "card-header" }));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.span, __VLS_intrinsicElements.span)({});
    var __VLS_280 = {}.ElSelect;
    /** @type {[typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, typeof __VLS_components.ElSelect, typeof __VLS_components.elSelect, ]} */ ;
    // @ts-ignore
    var __VLS_281 = __VLS_asFunctionalComponent(__VLS_280, new __VLS_280(__assign({ modelValue: (__VLS_ctx.selectedModel), placeholder: "选择AI模型", size: "small" }, { style: {} })));
    var __VLS_282 = __VLS_281.apply(void 0, __spreadArray([__assign({ modelValue: (__VLS_ctx.selectedModel), placeholder: "选择AI模型", size: "small" }, { style: {} })], __VLS_functionalComponentArgsRest(__VLS_281), false));
    __VLS_283.slots.default;
    var __VLS_284 = {}.ElOption;
    /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
    // @ts-ignore
    var __VLS_285 = __VLS_asFunctionalComponent(__VLS_284, new __VLS_284({
        label: "豆包",
        value: "doubao",
    }));
    var __VLS_286 = __VLS_285.apply(void 0, __spreadArray([{
            label: "豆包",
            value: "doubao",
        }], __VLS_functionalComponentArgsRest(__VLS_285), false));
    var __VLS_288 = {}.ElOption;
    /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
    // @ts-ignore
    var __VLS_289 = __VLS_asFunctionalComponent(__VLS_288, new __VLS_288({
        label: "DeepSeek",
        value: "deepseek",
    }));
    var __VLS_290 = __VLS_289.apply(void 0, __spreadArray([{
            label: "DeepSeek",
            value: "deepseek",
        }], __VLS_functionalComponentArgsRest(__VLS_289), false));
    var __VLS_292 = {}.ElOption;
    /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
    // @ts-ignore
    var __VLS_293 = __VLS_asFunctionalComponent(__VLS_292, new __VLS_292({
        label: "千问",
        value: "qianwen",
    }));
    var __VLS_294 = __VLS_293.apply(void 0, __spreadArray([{
            label: "千问",
            value: "qianwen",
        }], __VLS_functionalComponentArgsRest(__VLS_293), false));
    var __VLS_296 = {}.ElOption;
    /** @type {[typeof __VLS_components.ElOption, typeof __VLS_components.elOption, ]} */ ;
    // @ts-ignore
    var __VLS_297 = __VLS_asFunctionalComponent(__VLS_296, new __VLS_296({
        label: "元宝",
        value: "yuanbao",
    }));
    var __VLS_298 = __VLS_297.apply(void 0, __spreadArray([{
            label: "元宝",
            value: "yuanbao",
        }], __VLS_functionalComponentArgsRest(__VLS_297), false));
    var __VLS_283;
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "chat-container" }));
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "chat-messages" }, { ref: "chatMessages" }));
/** @type {typeof __VLS_ctx.chatMessages} */ ;
for (var _i = 0, _a = __VLS_getVForSourceType((__VLS_ctx.messages)); _i < _a.length; _i++) {
    var _b = _a[_i], message = _b[0], index = _b[1];
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ key: (index) }, { class: (['message-item', message.role]) }));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "message-avatar" }));
    if (message.role === 'user') {
        var __VLS_300 = {}.ElAvatar;
        /** @type {[typeof __VLS_components.ElAvatar, typeof __VLS_components.elAvatar, ]} */ ;
        // @ts-ignore
        var __VLS_301 = __VLS_asFunctionalComponent(__VLS_300, new __VLS_300(__assign({ icon: (__VLS_ctx.User) }, { class: "user-avatar" })));
        var __VLS_302 = __VLS_301.apply(void 0, __spreadArray([__assign({ icon: (__VLS_ctx.User) }, { class: "user-avatar" })], __VLS_functionalComponentArgsRest(__VLS_301), false));
    }
    else {
        var __VLS_304 = {}.ElAvatar;
        /** @type {[typeof __VLS_components.ElAvatar, typeof __VLS_components.elAvatar, ]} */ ;
        // @ts-ignore
        var __VLS_305 = __VLS_asFunctionalComponent(__VLS_304, new __VLS_304(__assign({ icon: (__VLS_ctx.ChatLineRound) }, { class: "ai-avatar" })));
        var __VLS_306 = __VLS_305.apply(void 0, __spreadArray([__assign({ icon: (__VLS_ctx.ChatLineRound) }, { class: "ai-avatar" })], __VLS_functionalComponentArgsRest(__VLS_305), false));
    }
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "message-content" }));
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "message-text" }));
    (message.content);
    __VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "message-time" }));
    (message.timestamp);
}
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "chat-input-area" }));
var __VLS_308 = {}.ElInput;
/** @type {[typeof __VLS_components.ElInput, typeof __VLS_components.elInput, ]} */ ;
// @ts-ignore
var __VLS_309 = __VLS_asFunctionalComponent(__VLS_308, new __VLS_308(__assign({ 'onKeyup': {} }, { modelValue: (__VLS_ctx.inputMessage), type: "textarea", placeholder: "请输入您的问题...", rows: (3), resize: "none" })));
var __VLS_310 = __VLS_309.apply(void 0, __spreadArray([__assign({ 'onKeyup': {} }, { modelValue: (__VLS_ctx.inputMessage), type: "textarea", placeholder: "请输入您的问题...", rows: (3), resize: "none" })], __VLS_functionalComponentArgsRest(__VLS_309), false));
var __VLS_312;
var __VLS_313;
var __VLS_314;
var __VLS_315 = {
    onKeyup: (__VLS_ctx.sendMessage)
};
var __VLS_311;
__VLS_asFunctionalElement(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)(__assign({ class: "chat-actions" }));
var __VLS_316 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_317 = __VLS_asFunctionalComponent(__VLS_316, new __VLS_316(__assign({ 'onClick': {} })));
var __VLS_318 = __VLS_317.apply(void 0, __spreadArray([__assign({ 'onClick': {} })], __VLS_functionalComponentArgsRest(__VLS_317), false));
var __VLS_320;
var __VLS_321;
var __VLS_322;
var __VLS_323 = {
    onClick: (__VLS_ctx.clearChat)
};
__VLS_319.slots.default;
var __VLS_319;
var __VLS_324 = {}.ElButton;
/** @type {[typeof __VLS_components.ElButton, typeof __VLS_components.elButton, typeof __VLS_components.ElButton, typeof __VLS_components.elButton, ]} */ ;
// @ts-ignore
var __VLS_325 = __VLS_asFunctionalComponent(__VLS_324, new __VLS_324(__assign({ 'onClick': {} }, { type: "primary", loading: (__VLS_ctx.sending) })));
var __VLS_326 = __VLS_325.apply(void 0, __spreadArray([__assign({ 'onClick': {} }, { type: "primary", loading: (__VLS_ctx.sending) })], __VLS_functionalComponentArgsRest(__VLS_325), false));
var __VLS_328;
var __VLS_329;
var __VLS_330;
var __VLS_331 = {
    onClick: (__VLS_ctx.sendMessage)
};
__VLS_327.slots.default;
var __VLS_327;
var __VLS_279;
/** @type {__VLS_StyleScopedClasses['ai-intelligence-view']} */ ;
/** @type {__VLS_StyleScopedClasses['page-header']} */ ;
/** @type {__VLS_StyleScopedClasses['header-content']} */ ;
/** @type {__VLS_StyleScopedClasses['config-card']} */ ;
/** @type {__VLS_StyleScopedClasses['card-header']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-tabs']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-form']} */ ;
/** @type {__VLS_StyleScopedClasses['slider-value']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-form']} */ ;
/** @type {__VLS_StyleScopedClasses['slider-value']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-form']} */ ;
/** @type {__VLS_StyleScopedClasses['slider-value']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-form']} */ ;
/** @type {__VLS_StyleScopedClasses['slider-value']} */ ;
/** @type {__VLS_StyleScopedClasses['chat-card']} */ ;
/** @type {__VLS_StyleScopedClasses['card-header']} */ ;
/** @type {__VLS_StyleScopedClasses['chat-container']} */ ;
/** @type {__VLS_StyleScopedClasses['chat-messages']} */ ;
/** @type {__VLS_StyleScopedClasses['message-avatar']} */ ;
/** @type {__VLS_StyleScopedClasses['user-avatar']} */ ;
/** @type {__VLS_StyleScopedClasses['ai-avatar']} */ ;
/** @type {__VLS_StyleScopedClasses['message-content']} */ ;
/** @type {__VLS_StyleScopedClasses['message-text']} */ ;
/** @type {__VLS_StyleScopedClasses['message-time']} */ ;
/** @type {__VLS_StyleScopedClasses['chat-input-area']} */ ;
/** @type {__VLS_StyleScopedClasses['chat-actions']} */ ;
// @ts-ignore
var __VLS_21 = __VLS_20, __VLS_87 = __VLS_86, __VLS_153 = __VLS_152, __VLS_219 = __VLS_218;
var __VLS_dollars;
var __VLS_self = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {
            User: icons_vue_1.User,
            ChatLineRound: icons_vue_1.ChatLineRound,
            doubaoForm: doubaoForm,
            deepseekForm: deepseekForm,
            qianwenForm: qianwenForm,
            yuanbaoForm: yuanbaoForm,
            chatMessages: chatMessages,
            activeTab: activeTab,
            sending: sending,
            selectedModel: selectedModel,
            messages: messages,
            inputMessage: inputMessage,
            aiConfig: aiConfig,
            saveConfig: saveConfig,
            testConnection: testConnection,
            sendMessage: sendMessage,
            clearChat: clearChat,
        };
    },
});
exports.default = (await Promise.resolve().then(function () { return require('vue'); })).defineComponent({
    setup: function () {
        return {};
    },
});
; /* PartiallyEnd: #4569/main.vue */
;
