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
    const asciiInput = document.getElementById(isEncrypt ? 'plaintextASCII' : 'ciphertextASCII');
    const label = binaryInput.previousElementSibling;

    if (mode === 'ascii') {
        binaryInput.style.display = 'none';
        binaryInput.value = '';
        asciiInput.style.display = 'block';
        asciiInput.value = '';
        asciiInput.focus();
        label.textContent = isEncrypt ? '明文 (ASCII 字符):' : '密文 (ASCII 字符):';
    } else {
        asciiInput.style.display = 'none';
        asciiInput.value = '';
        binaryInput.style.display = 'block';
        binaryInput.value = '';
        binaryInput.focus();
        label.textContent = isEncrypt ? '明文 (8位二进制):' : '密文 (8位二进制):';
    }
}

function buildAsciiMessage(label, content) {
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

document.getElementById('encryptForm').addEventListener('submit', async (e) => {
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
            if (data.ciphertext_ascii) {
                showResult('encryptResult', buildAsciiMessage('ASCII 密文', data.ciphertext_ascii), true);
                document.getElementById('ciphertextASCII').value = data.ciphertext_ascii;
                document.querySelector('input[name="decryptMode"][value="ascii"]').checked = true;
                toggleInputMode('decryptForm', 'ascii');
            } else {
                showResult('encryptResult', `密文: ${data.ciphertext}`, true);
                document.getElementById('ciphertext').value = data.ciphertext;
                document.querySelector('input[name="decryptMode"][value="binary"]').checked = true;
                toggleInputMode('decryptForm', 'binary');
            }
            document.getElementById('decryptKey').value = keyValue;
        } else {
            showResult('encryptResult', `错误: ${data.message}`, false);
        }
    } catch (error) {
        showResult('encryptResult', `网络错误: ${error.message}`, false);
    }
});

document.getElementById('decryptForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const ciphertext = document.getElementById('ciphertext');
    const ciphertextASCII = document.getElementById('ciphertextASCII');
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
        const asciiValue = ciphertextASCII.value;
        if (!asciiValue) {
            showResult('decryptResult', '密文错误: ASCII 文本不能为空', false);
            return;
        }
        payload.ciphertext_ascii = asciiValue;
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
                showResult('decryptResult', buildAsciiMessage('ASCII 明文', data.plaintext_ascii), true);
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

document.addEventListener('DOMContentLoaded', () => {
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

