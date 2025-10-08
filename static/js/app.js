// 验证二进制输入
function validateBinaryInput(value, length) {
    if (value.length !== length) {
        return `输入必须是${length}位`;
    }
    if (!/^[01]+$/.test(value)) {
        return '只能包含0和1';
    }
    return null;
}

function toggleInputMode(formId, mode) {
    const isEncrypt = formId === 'encryptForm';
    const binaryInput = document.getElementById(isEncrypt ? 'plaintext' : 'ciphertext');
    const altInput = document.getElementById(isEncrypt ? 'plaintextASCII' : 'ciphertextBase64');
    const label = binaryInput.previousElementSibling;

    if (mode === 'ascii') {
        binaryInput.style.display = 'none';
        binaryInput.value = '';
        altInput.style.display = 'block';
        altInput.value = '';
        altInput.focus();
        label.textContent = isEncrypt ? '明文 (ASCII 字符):' : '密文 (Base64):';
    } else {
        altInput.style.display = 'none';
        altInput.value = '';
        binaryInput.style.display = 'block';
        binaryInput.value = '';
        binaryInput.focus();
        label.textContent = isEncrypt ? '明文 (8位二进制):' : '密文 (8位二进制):';
    }
}

function buildTextResult(label, content) {
    const container = document.createElement('div');

    const createBlock = (title, text, type, lengthMeta) => {
        const block = document.createElement('div');
        block.className = 'result-block';

        const titleEl = document.createElement('strong');
        titleEl.textContent = title;
        block.appendChild(titleEl);

        const pre = document.createElement('pre');
        pre.textContent = text === '' ? '(空)' : text;
        block.appendChild(pre);

        const actions = document.createElement('div');
        actions.className = 'result-actions';

        const copyBtn = document.createElement('button');
        copyBtn.type = 'button';
        copyBtn.className = 'btn-copy';
        copyBtn.textContent = '复制';
        copyBtn.dataset.copy = type;
        actions.appendChild(copyBtn);

        block.appendChild(actions);

        const meta = document.createElement('div');
        meta.className = 'result-meta';
        if (typeof lengthMeta === 'number') {
            const lengthSpan = document.createElement('span');
            lengthSpan.textContent = `长度: ${lengthMeta}`;
            meta.appendChild(lengthSpan);
        } else {
            meta.appendChild(document.createElement('span'));
        }

        const indicator = document.createElement('span');
        indicator.className = 'copy-indicator';
        indicator.style.display = 'none';
        indicator.textContent = '已复制';
        meta.appendChild(indicator);

        block.appendChild(meta);

        return { block, copyBtn, indicator };
    };

    const rawBlock = createBlock(label, content, 'raw', content.length);

    container.appendChild(rawBlock.block);

    const copyMap = {
        raw: content,
    };

    [rawBlock].forEach(({ copyBtn, indicator }) => {
        copyBtn.addEventListener('click', async () => {
            const type = copyBtn.dataset.copy;
            const textToCopy = copyMap[type] ?? '';
            try {
                await navigator.clipboard.writeText(textToCopy);
                indicator.style.display = 'inline';
                setTimeout(() => {
                    indicator.style.display = 'none';
                }, 1500);
            } catch (err) {
                console.error('复制失败:', err);
            }
        });
    });

    return container;
}

// 构建暴力破解结果显示
function buildBruteForceResult(data) {
    const container = document.createElement('div');
    container.className = 'brute-force-results';

    // 成功消息
    const message = document.createElement('div');
    message.className = 'result-message';
    message.innerHTML = `<strong>${data.message}</strong>`;
    container.appendChild(message);

    if (data.keys && data.keys.length > 0) {
        // 密钥列表
        const keysList = document.createElement('div');
        keysList.className = 'keys-list';
        
        const listTitle = document.createElement('div');
        listTitle.innerHTML = '<strong>找到的密钥：</strong>';
        listTitle.style.marginBottom = '8px';
        keysList.appendChild(listTitle);

        data.keys.forEach((key, index) => {
            const keyItem = document.createElement('div');
            keyItem.className = 'key-item';
            
            const keyBinary = document.createElement('span');
            keyBinary.className = 'key-binary';
            keyBinary.textContent = key;
            
            const keyDecimal = document.createElement('span');
            keyDecimal.className = 'key-decimal';
            keyDecimal.textContent = `十进制: ${data.keys_decimal[index]}`;
            
            keyItem.appendChild(keyBinary);
            keyItem.appendChild(keyDecimal);
            keysList.appendChild(keyItem);
        });

        container.appendChild(keysList);

        // 统计信息
        const statsInfo = document.createElement('div');
        statsInfo.className = 'stats-info';
        
        const countSpan = document.createElement('span');
        countSpan.className = 'stats-count';
        countSpan.textContent = `找到 ${data.key_count} 个密钥`;
        
        const timeSpan = document.createElement('span');
        timeSpan.className = 'stats-time';
        timeSpan.textContent = `耗时: ${data.time}`;
        
        statsInfo.appendChild(countSpan);
        statsInfo.appendChild(timeSpan);
        container.appendChild(statsInfo);
    }

    return container;
}

function showResult(elementId, content, isSuccess = true) {
    const resultElement = document.getElementById(elementId);
    resultElement.innerHTML = '';
    if (content instanceof HTMLElement) {
        resultElement.appendChild(content);
    } else {
        resultElement.textContent = content;
    }
    resultElement.className = `result ${isSuccess ? 'success' : 'error'}`;
    resultElement.style.display = 'flex';
}

function showLoading(elementId) {
    const resultElement = document.getElementById(elementId);
    resultElement.innerHTML = '<div class="loading"></div>处理中...';
    resultElement.className = 'result';
    resultElement.style.display = 'flex';
}

// 将在DOMContentLoaded中绑定
function bindEncryptForm() {
    const form = document.getElementById('encryptForm');
    if (!form) return;
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const plaintext = document.getElementById('plaintext');
        const plaintextASCII = document.getElementById('plaintextASCII');
        const key = document.getElementById('encryptKey');
        const mode = document.querySelector('input[name="encryptMode"]:checked').value;

        const keyValue = key.value.trim();
        const keyError = validateBinaryInput(keyValue, 10);
        if (keyError) {
            showResult('encryptResult', `密钥错误: ${keyError}`, false);
            return;
        }

        const payload = { key: keyValue };

        if (mode === 'ascii') {
            const asciiValue = plaintextASCII.value;
            if (!asciiValue) {
                showResult('encryptResult', '明文错误: ASCII 文本不能为空', false);
                return;
            }
            payload.plaintext_ascii = asciiValue;
        } else {
            const binaryValue = plaintext.value.trim();
            const plaintextError = validateBinaryInput(binaryValue, 8);
            if (plaintextError) {
                showResult('encryptResult', `明文错误: ${plaintextError}`, false);
                return;
            }
            payload.plaintext = binaryValue;
        }

        showLoading('encryptResult');

        try {
            const response = await fetch('/api/encrypt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload)
            });

            const data = await response.json();

            if (data.success) {
                if (data.ciphertext_base64) {
                    showResult('encryptResult', buildTextResult('Base64 密文', data.ciphertext_base64), true);
                } else {
                    const binary = data.ciphertext_binary || data.ciphertext;
                    showResult('encryptResult', `密文: ${binary ?? '未知'}`, true);
                }
            } else {
                showResult('encryptResult', `错误: ${data.message}`, false);
            }
        } catch (error) {
            showResult('encryptResult', `网络错误: ${error.message}`, false);
        }
    });
}

// 将在DOMContentLoaded中绑定
function bindDecryptForm() {
    const form = document.getElementById('decryptForm');
    if (!form) return;
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const ciphertext = document.getElementById('ciphertext');
        const ciphertextBase64 = document.getElementById('ciphertextBase64');
        const key = document.getElementById('decryptKey');
        const mode = document.querySelector('input[name="decryptMode"]:checked').value;

        const keyValue = key.value.trim();
        const keyError = validateBinaryInput(keyValue, 10);
        if (keyError) {
            showResult('decryptResult', `密钥错误: ${keyError}`, false);
            return;
        }

        const payload = { key: keyValue };

        if (mode === 'ascii') {
            const base64Value = ciphertextBase64.value.trim();
            if (!base64Value) {
                showResult('decryptResult', '密文错误: Base64 文本不能为空', false);
                return;
            }
            payload.ciphertext_base64 = base64Value;
        } else {
            const binaryValue = ciphertext.value.trim();
            const ciphertextError = validateBinaryInput(binaryValue, 8);
            if (ciphertextError) {
                showResult('decryptResult', `密文错误: ${ciphertextError}`, false);
                return;
            }
            payload.ciphertext = binaryValue;
        }

        showLoading('decryptResult');

        try {
            const response = await fetch('/api/decrypt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload)
            });

            const data = await response.json();

            if (data.success) {
                if (data.plaintext_ascii) {
                    showResult('decryptResult', buildTextResult('ASCII 明文', data.plaintext_ascii), true);
                } else {
                    showResult('decryptResult', `明文: ${data.plaintext}`, true);
                }
            } else {
                showResult('decryptResult', `错误: ${data.message}`, false);
            }
        } catch (error) {
            showResult('decryptResult', `网络错误: ${error.message}`, false);
        }
    });
}

// 将在DOMContentLoaded中绑定
function bindBruteForceForm() {
    const form = document.getElementById('bruteForceForm');
    if (!form) return;
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const plaintext = document.getElementById('brutePlaintext');
        const ciphertext = document.getElementById('bruteCiphertext');

        const plaintextValue = plaintext.value.trim();
        const ciphertextValue = ciphertext.value.trim();

        // 验证输入
        const plaintextError = validateBinaryInput(plaintextValue, 8);
        if (plaintextError) {
            showResult('bruteForceResult', `明文错误: ${plaintextError}`, false);
            return;
        }

        const ciphertextError = validateBinaryInput(ciphertextValue, 8);
        if (ciphertextError) {
            showResult('bruteForceResult', `密文错误: ${ciphertextError}`, false);
            return;
        }

        const payload = {
            plaintext: plaintextValue,
            ciphertext: ciphertextValue
        };

        showLoading('bruteForceResult');

        try {
            const response = await fetch('/api/blasting', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(payload)
            });

            const data = await response.json();

            if (data.success) {
                showResult('bruteForceResult', buildBruteForceResult(data), true);
            } else {
                showResult('bruteForceResult', `${data.message} (耗时: ${data.time || '未知'})`, false);
            }
        } catch (error) {
            showResult('bruteForceResult', `网络错误: ${error.message}`, false);
        }
    });
}

function bindBinaryInputSanitizer() {
    document.querySelectorAll('input[pattern]').forEach(input => {
        input.addEventListener('input', (e) => {
            e.target.value = e.target.value.replace(/[^01]/g, '');
        });

        input.addEventListener('paste', (e) => {
            e.preventDefault();
            const paste = (e.clipboardData || window.clipboardData).getData('text');
            const cleanPaste = paste.replace(/[^01]/g, '');
            e.target.value = cleanPaste.substring(0, e.target.maxLength);
        });
    });
}

document.addEventListener('DOMContentLoaded', () => {
    // 绑定表单事件
    bindEncryptForm();
    bindDecryptForm();
    bindBruteForceForm();
    bindBinaryInputSanitizer();
    
    // 绑定模式切换事件
    document.querySelectorAll('input[name="encryptMode"]').forEach(radio => {
        radio.addEventListener('change', (event) => {
            toggleInputMode('encryptForm', event.target.value);
        });
    });

    document.querySelectorAll('input[name="decryptMode"]').forEach(radio => {
        radio.addEventListener('change', (event) => {
            toggleInputMode('decryptForm', event.target.value);
        });
    });
});

