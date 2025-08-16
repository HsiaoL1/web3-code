# 第八章：代币标准与 DeFi

## 本章学习目标

- 深入理解 ERC-20 代币标准及其实现
- 掌握 ERC-721 NFT 标准的核心概念
- 学习 ERC-1155 多代币标准
- 了解 DeFi 协议的基本原理
- 实现代币交换和流动性挖矿
- 掌握代币安全最佳实践

## 8.1 ERC-20 代币标准

ERC-20 是以太坊上最重要的代币标准，定义了同质化代币的接口规范。

### 8.1.1 ERC-20 标准接口

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// ERC-20 标准接口
interface IERC20 {
    // 返回代币总供应量
    function totalSupply() external view returns (uint256);

    // 返回指定账户的代币余额
    function balanceOf(address account) external view returns (uint256);

    // 转移代币到指定账户
    function transfer(address to, uint256 amount) external returns (bool);

    // 返回 spender 被允许从 owner 账户提取的代币数量
    function allowance(address owner, address spender) external view returns (uint256);

    // 批准 spender 从调用者账户提取指定数量的代币
    function approve(address spender, uint256 amount) external returns (bool);

    // 从 from 账户转移代币到 to 账户（需要事先获得授权）
    function transferFrom(address from, address to, uint256 amount) external returns (bool);

    // 代币转移事件
    event Transfer(address indexed from, address indexed to, uint256 value);

    // 授权事件
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

// ERC-20 元数据扩展接口
interface IERC20Metadata is IERC20 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function decimals() external view returns (uint8);
}
```

### 8.1.2 基础 ERC-20 实现

```solidity
contract BasicERC20 is IERC20, IERC20Metadata {
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;

    uint256 private _totalSupply;
    string private _name;
    string private _symbol;
    uint8 private _decimals;

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        uint256 totalSupply_
    ) {
        _name = name_;
        _symbol = symbol_;
        _decimals = decimals_;
        _totalSupply = totalSupply_ * 10**decimals_;
        _balances[msg.sender] = _totalSupply;
        emit Transfer(address(0), msg.sender, _totalSupply);
    }

    function name() public view override returns (string memory) {
        return _name;
    }

    function symbol() public view override returns (string memory) {
        return _symbol;
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    function totalSupply() public view override returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view override returns (uint256) {
        return _balances[account];
    }

    function transfer(address to, uint256 amount) public override returns (bool) {
        address owner = msg.sender;
        _transfer(owner, to, amount);
        return true;
    }

    function allowance(address owner, address spender) public view override returns (uint256) {
        return _allowances[owner][spender];
    }

    function approve(address spender, uint256 amount) public override returns (bool) {
        address owner = msg.sender;
        _approve(owner, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) public override returns (bool) {
        address spender = msg.sender;
        _spendAllowance(from, spender, amount);
        _transfer(from, to, amount);
        return true;
    }

    function _transfer(address from, address to, uint256 amount) internal {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");

        uint256 fromBalance = _balances[from];
        require(fromBalance >= amount, "ERC20: transfer amount exceeds balance");

        unchecked {
            _balances[from] = fromBalance - amount;
            _balances[to] += amount;
        }

        emit Transfer(from, to, amount);
    }

    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    function _spendAllowance(address owner, address spender, uint256 amount) internal {
        uint256 currentAllowance = allowance(owner, spender);
        if (currentAllowance != type(uint256).max) {
            require(currentAllowance >= amount, "ERC20: insufficient allowance");
            unchecked {
                _approve(owner, spender, currentAllowance - amount);
            }
        }
    }
}
```

### 8.1.3 高级 ERC-20 功能

```solidity
contract AdvancedERC20 is BasicERC20 {
    address public owner;
    bool public paused = false;
    mapping(address => bool) public blacklisted;

    // 代币销毁事件
    event Burn(address indexed from, uint256 value);
    // 代币铸造事件
    event Mint(address indexed to, uint256 value);
    // 暂停状态变更事件
    event PauseChanged(bool paused);
    // 黑名单状态变更事件
    event BlacklistChanged(address indexed account, bool blacklisted);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    modifier notBlacklisted(address account) {
        require(!blacklisted[account], "Account is blacklisted");
        _;
    }

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        uint256 totalSupply_
    ) BasicERC20(name_, symbol_, decimals_, totalSupply_) {
        owner = msg.sender;
    }

    // 铸造代币
    function mint(address to, uint256 amount) public onlyOwner {
        require(to != address(0), "ERC20: mint to the zero address");

        _totalSupply += amount;
        unchecked {
            _balances[to] += amount;
        }

        emit Transfer(address(0), to, amount);
        emit Mint(to, amount);
    }

    // 销毁代币
    function burn(uint256 amount) public {
        address account = msg.sender;
        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "ERC20: burn amount exceeds balance");

        unchecked {
            _balances[account] = accountBalance - amount;
            _totalSupply -= amount;
        }

        emit Transfer(account, address(0), amount);
        emit Burn(account, amount);
    }

    // 从指定账户销毁代币
    function burnFrom(address account, uint256 amount) public {
        _spendAllowance(account, msg.sender, amount);

        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "ERC20: burn amount exceeds balance");

        unchecked {
            _balances[account] = accountBalance - amount;
            _totalSupply -= amount;
        }

        emit Transfer(account, address(0), amount);
        emit Burn(account, amount);
    }

    // 暂停/恢复合约
    function setPaused(bool _paused) public onlyOwner {
        paused = _paused;
        emit PauseChanged(_paused);
    }

    // 黑名单管理
    function setBlacklisted(address account, bool _blacklisted) public onlyOwner {
        blacklisted[account] = _blacklisted;
        emit BlacklistChanged(account, _blacklisted);
    }

    // 重写转账函数以添加限制
    function _transfer(address from, address to, uint256 amount)
        internal
        override
        whenNotPaused
        notBlacklisted(from)
        notBlacklisted(to)
    {
        super._transfer(from, to, amount);
    }

    // 转移所有权
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner cannot be zero address");
        owner = newOwner;
    }

    // 增加授权额度
    function increaseAllowance(address spender, uint256 addedValue) public returns (bool) {
        address owner = msg.sender;
        _approve(owner, spender, allowance(owner, spender) + addedValue);
        return true;
    }

    // 减少授权额度
    function decreaseAllowance(address spender, uint256 subtractedValue) public returns (bool) {
        address owner = msg.sender;
        uint256 currentAllowance = allowance(owner, spender);
        require(currentAllowance >= subtractedValue, "ERC20: decreased allowance below zero");
        unchecked {
            _approve(owner, spender, currentAllowance - subtractedValue);
        }
        return true;
    }
}
```

## 8.2 ERC-721 NFT 标准

ERC-721 定义了非同质化代币（NFT）的标准接口。

### 8.2.1 ERC-721 标准接口

```solidity
interface IERC721 {
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);

    function balanceOf(address owner) external view returns (uint256 balance);
    function ownerOf(uint256 tokenId) external view returns (address owner);
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external;
    function safeTransferFrom(address from, address to, uint256 tokenId) external;
    function transferFrom(address from, address to, uint256 tokenId) external;
    function approve(address to, uint256 tokenId) external;
    function setApprovalForAll(address operator, bool approved) external;
    function getApproved(uint256 tokenId) external view returns (address operator);
    function isApprovedForAll(address owner, address operator) external view returns (bool);
}

interface IERC721Metadata is IERC721 {
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function tokenURI(uint256 tokenId) external view returns (string memory);
}

// ERC-721 接收器接口
interface IERC721Receiver {
    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external returns (bytes4);
}
```

### 8.2.2 基础 NFT 实现

```solidity
contract BasicNFT is IERC721, IERC721Metadata {
    using Strings for uint256;

    string private _name;
    string private _symbol;
    string private _baseTokenURI;

    mapping(uint256 => address) private _owners;
    mapping(address => uint256) private _balances;
    mapping(uint256 => address) private _tokenApprovals;
    mapping(address => mapping(address => bool)) private _operatorApprovals;

    uint256 private _currentTokenId = 0;

    constructor(string memory name_, string memory symbol_, string memory baseURI_) {
        _name = name_;
        _symbol = symbol_;
        _baseTokenURI = baseURI_;
    }

    function name() public view override returns (string memory) {
        return _name;
    }

    function symbol() public view override returns (string memory) {
        return _symbol;
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        require(_exists(tokenId), "ERC721: URI query for nonexistent token");
        return bytes(_baseTokenURI).length > 0
            ? string(abi.encodePacked(_baseTokenURI, tokenId.toString()))
            : "";
    }

    function balanceOf(address owner) public view override returns (uint256) {
        require(owner != address(0), "ERC721: balance query for the zero address");
        return _balances[owner];
    }

    function ownerOf(uint256 tokenId) public view override returns (address) {
        address owner = _owners[tokenId];
        require(owner != address(0), "ERC721: owner query for nonexistent token");
        return owner;
    }

    function approve(address to, uint256 tokenId) public override {
        address owner = ownerOf(tokenId);
        require(to != owner, "ERC721: approval to current owner");
        require(
            msg.sender == owner || isApprovedForAll(owner, msg.sender),
            "ERC721: approve caller is not owner nor approved for all"
        );

        _approve(to, tokenId);
    }

    function getApproved(uint256 tokenId) public view override returns (address) {
        require(_exists(tokenId), "ERC721: approved query for nonexistent token");
        return _tokenApprovals[tokenId];
    }

    function setApprovalForAll(address operator, bool approved) public override {
        require(operator != msg.sender, "ERC721: approve to caller");
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function isApprovedForAll(address owner, address operator) public view override returns (bool) {
        return _operatorApprovals[owner][operator];
    }

    function transferFrom(address from, address to, uint256 tokenId) public override {
        require(_isApprovedOrOwner(msg.sender, tokenId), "ERC721: transfer caller is not owner nor approved");
        _transfer(from, to, tokenId);
    }

    function safeTransferFrom(address from, address to, uint256 tokenId) public override {
        safeTransferFrom(from, to, tokenId, "");
    }

    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public override {
        require(_isApprovedOrOwner(msg.sender, tokenId), "ERC721: transfer caller is not owner nor approved");
        _safeTransfer(from, to, tokenId, data);
    }

    function _exists(uint256 tokenId) internal view returns (bool) {
        return _owners[tokenId] != address(0);
    }

    function _isApprovedOrOwner(address spender, uint256 tokenId) internal view returns (bool) {
        require(_exists(tokenId), "ERC721: operator query for nonexistent token");
        address owner = ownerOf(tokenId);
        return (spender == owner || getApproved(tokenId) == spender || isApprovedForAll(owner, spender));
    }

    function _mint(address to, uint256 tokenId) internal {
        require(to != address(0), "ERC721: mint to the zero address");
        require(!_exists(tokenId), "ERC721: token already minted");

        _balances[to] += 1;
        _owners[tokenId] = to;

        emit Transfer(address(0), to, tokenId);
    }

    function _transfer(address from, address to, uint256 tokenId) internal {
        require(ownerOf(tokenId) == from, "ERC721: transfer from incorrect owner");
        require(to != address(0), "ERC721: transfer to the zero address");

        _approve(address(0), tokenId);

        _balances[from] -= 1;
        _balances[to] += 1;
        _owners[tokenId] = to;

        emit Transfer(from, to, tokenId);
    }

    function _approve(address to, uint256 tokenId) internal {
        _tokenApprovals[tokenId] = to;
        emit Approval(ownerOf(tokenId), to, tokenId);
    }

    function _safeTransfer(address from, address to, uint256 tokenId, bytes memory data) internal {
        _transfer(from, to, tokenId);
        require(_checkOnERC721Received(from, to, tokenId, data), "ERC721: transfer to non ERC721Receiver implementer");
    }

    function _checkOnERC721Received(address from, address to, uint256 tokenId, bytes memory data) private returns (bool) {
        if (to.code.length > 0) {
            try IERC721Receiver(to).onERC721Received(msg.sender, from, tokenId, data) returns (bytes4 retval) {
                return retval == IERC721Receiver.onERC721Received.selector;
            } catch (bytes memory reason) {
                if (reason.length == 0) {
                    revert("ERC721: transfer to non ERC721Receiver implementer");
                } else {
                    assembly {
                        revert(add(32, reason), mload(reason))
                    }
                }
            }
        } else {
            return true;
        }
    }

    // 公共铸造函数
    function mint(address to) public returns (uint256) {
        uint256 tokenId = _currentTokenId++;
        _mint(to, tokenId);
        return tokenId;
    }
}
```

### 8.2.3 高级 NFT 功能

```solidity
contract AdvancedNFT is BasicNFT {
    struct TokenInfo {
        string name;
        string description;
        string image;
        mapping(string => string) attributes;
        string[] attributeKeys;
    }

    mapping(uint256 => TokenInfo) private _tokenInfo;
    mapping(address => bool) public minters;
    uint256 public maxSupply;
    uint256 public mintPrice;
    bool public publicMintEnabled = false;

    event TokenInfoUpdated(uint256 indexed tokenId);
    event MinterAdded(address indexed minter);
    event MinterRemoved(address indexed minter);

    modifier onlyMinter() {
        require(minters[msg.sender], "Caller is not a minter");
        _;
    }

    constructor(
        string memory name_,
        string memory symbol_,
        string memory baseURI_,
        uint256 maxSupply_,
        uint256 mintPrice_
    ) BasicNFT(name_, symbol_, baseURI_) {
        maxSupply = maxSupply_;
        mintPrice = mintPrice_;
        minters[msg.sender] = true;
    }

    function addMinter(address minter) external onlyMinter {
        minters[minter] = true;
        emit MinterAdded(minter);
    }

    function removeMinter(address minter) external onlyMinter {
        minters[minter] = false;
        emit MinterRemoved(minter);
    }

    function setPublicMintEnabled(bool enabled) external onlyMinter {
        publicMintEnabled = enabled;
    }

    function setMintPrice(uint256 price) external onlyMinter {
        mintPrice = price;
    }

    function publicMint() external payable returns (uint256) {
        require(publicMintEnabled, "Public minting is not enabled");
        require(msg.value >= mintPrice, "Insufficient payment");
        require(_currentTokenId < maxSupply, "Max supply reached");

        return mint(msg.sender);
    }

    function batchMint(address to, uint256 quantity) external onlyMinter returns (uint256[] memory) {
        require(_currentTokenId + quantity <= maxSupply, "Would exceed max supply");

        uint256[] memory tokenIds = new uint256[](quantity);
        for (uint256 i = 0; i < quantity; i++) {
            tokenIds[i] = mint(to);
        }
        return tokenIds;
    }

    function setTokenInfo(
        uint256 tokenId,
        string memory name_,
        string memory description_,
        string memory image_
    ) external onlyMinter {
        require(_exists(tokenId), "Token does not exist");

        TokenInfo storage info = _tokenInfo[tokenId];
        info.name = name_;
        info.description = description_;
        info.image = image_;

        emit TokenInfoUpdated(tokenId);
    }

    function setTokenAttribute(
        uint256 tokenId,
        string memory key,
        string memory value
    ) external onlyMinter {
        require(_exists(tokenId), "Token does not exist");

        TokenInfo storage info = _tokenInfo[tokenId];

        // 如果是新属性，添加到键数组
        if (bytes(info.attributes[key]).length == 0) {
            info.attributeKeys.push(key);
        }

        info.attributes[key] = value;
        emit TokenInfoUpdated(tokenId);
    }

    function getTokenInfo(uint256 tokenId)
        external
        view
        returns (
            string memory name_,
            string memory description_,
            string memory image_,
            string[] memory attributeKeys,
            string[] memory attributeValues
        )
    {
        require(_exists(tokenId), "Token does not exist");

        TokenInfo storage info = _tokenInfo[tokenId];

        attributeKeys = info.attributeKeys;
        attributeValues = new string[](attributeKeys.length);

        for (uint256 i = 0; i < attributeKeys.length; i++) {
            attributeValues[i] = info.attributes[attributeKeys[i]];
        }

        return (info.name, info.description, info.image, attributeKeys, attributeValues);
    }

    function withdraw() external onlyMinter {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        (bool success, ) = payable(msg.sender).call{value: balance}("");
        require(success, "Withdrawal failed");
    }

    // 覆盖 tokenURI 以支持自定义元数据
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        require(_exists(tokenId), "ERC721: URI query for nonexistent token");

        TokenInfo storage info = _tokenInfo[tokenId];

        if (bytes(info.name).length > 0) {
            // 返回自定义元数据 JSON
            return _buildTokenURI(tokenId, info);
        } else {
            // 回退到基础实现
            return super.tokenURI(tokenId);
        }
    }

    function _buildTokenURI(uint256 tokenId, TokenInfo storage info)
        private
        view
        returns (string memory)
    {
        // 构建 JSON 元数据
        string memory attributes = _buildAttributes(info);

        return string(abi.encodePacked(
            'data:application/json;base64,',
            Base64.encode(bytes(abi.encodePacked(
                '{"name":"', info.name, '",',
                '"description":"', info.description, '",',
                '"image":"', info.image, '",',
                '"attributes":[', attributes, ']}'
            )))
        ));
    }

    function _buildAttributes(TokenInfo storage info) private view returns (string memory) {
        if (info.attributeKeys.length == 0) {
            return "";
        }

        string memory attributes = "";
        for (uint256 i = 0; i < info.attributeKeys.length; i++) {
            if (i > 0) {
                attributes = string(abi.encodePacked(attributes, ","));
            }
            attributes = string(abi.encodePacked(
                attributes,
                '{"trait_type":"', info.attributeKeys[i], '",',
                '"value":"', info.attributes[info.attributeKeys[i]], '"}'
            ));
        }

        return attributes;
    }
}

// Base64 编码库
library Base64 {
    bytes internal constant TABLE = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

    function encode(bytes memory data) internal pure returns (string memory) {
        uint256 len = data.length;
        if (len == 0) return "";

        uint256 encodedLen = 4 * ((len + 2) / 3);

        bytes memory result = new bytes(encodedLen + 32);

        bytes memory table = TABLE;

        assembly {
            let tablePtr := add(table, 1)
            let resultPtr := add(result, 32)

            for {
                let i := 0
            } lt(i, len) {

            } {
                i := add(i, 3)
                let input := and(mload(add(data, i)), 0xffffff)

                let out := mload(add(tablePtr, and(shr(18, input), 0x3F)))
                out := shl(8, out)
                out := add(out, and(mload(add(tablePtr, and(shr(12, input), 0x3F))), 0xFF))
                out := shl(8, out)
                out := add(out, and(mload(add(tablePtr, and(shr(6, input), 0x3F))), 0xFF))
                out := shl(8, out)
                out := add(out, and(mload(add(tablePtr, and(input, 0x3F))), 0xFF))
                out := shl(224, out)

                mstore(resultPtr, out)

                resultPtr := add(resultPtr, 4)
            }

            switch mod(len, 3)
            case 1 {
                mstore(sub(resultPtr, 2), shl(240, 0x3d3d))
            }
            case 2 {
                mstore(sub(resultPtr, 1), shl(248, 0x3d))
            }

            mstore(result, encodedLen)
        }

        return string(result);
    }
}

// String 工具库
library Strings {
    bytes16 private constant _HEX_SYMBOLS = "0123456789abcdef";

    function toString(uint256 value) internal pure returns (string memory) {
        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }
}
```

## 8.3 ERC-1155 多代币标准

ERC-1155 支持在单个合约中管理多种代币类型。

### 8.3.1 ERC-1155 接口

```solidity
interface IERC1155 {
    event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value);
    event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values);
    event ApprovalForAll(address indexed account, address indexed operator, bool approved);
    event URI(string value, uint256 indexed id);

    function balanceOf(address account, uint256 id) external view returns (uint256);
    function balanceOfBatch(address[] calldata accounts, uint256[] calldata ids) external view returns (uint256[] memory);
    function setApprovalForAll(address operator, bool approved) external;
    function isApprovedForAll(address account, address operator) external view returns (bool);
    function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes calldata data) external;
    function safeBatchTransferFrom(address from, address to, uint256[] calldata ids, uint256[] calldata amounts, bytes calldata data) external;
}

interface IERC1155MetadataURI is IERC1155 {
    function uri(uint256 id) external view returns (string memory);
}

interface IERC1155Receiver {
    function onERC1155Received(
        address operator,
        address from,
        uint256 id,
        uint256 value,
        bytes calldata data
    ) external returns (bytes4);

    function onERC1155BatchReceived(
        address operator,
        address from,
        uint256[] calldata ids,
        uint256[] calldata values,
        bytes calldata data
    ) external returns (bytes4);
}
```

### 8.3.2 ERC-1155 实现

```solidity
contract GameItems is IERC1155, IERC1155MetadataURI {
    mapping(uint256 => mapping(address => uint256)) private _balances;
    mapping(address => mapping(address => bool)) private _operatorApprovals;
    mapping(uint256 => string) private _tokenURIs;

    address public owner;
    string private _uri;

    // 游戏物品类型
    uint256 public constant SWORD = 1;
    uint256 public constant SHIELD = 2;
    uint256 public constant POTION = 3;
    uint256 public constant ARMOR = 4;
    uint256 public constant BOW = 5;

    constructor(string memory uri_) {
        _uri = uri_;
        owner = msg.sender;

        // 初始化一些物品
        _mint(msg.sender, SWORD, 100, "");
        _mint(msg.sender, SHIELD, 50, "");
        _mint(msg.sender, POTION, 1000, "");
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function uri(uint256 id) public view override returns (string memory) {
        string memory tokenURI = _tokenURIs[id];

        if (bytes(tokenURI).length > 0) {
            return tokenURI;
        }

        return string(abi.encodePacked(_uri, Strings.toString(id)));
    }

    function setURI(string memory newuri) external onlyOwner {
        _uri = newuri;
    }

    function setTokenURI(uint256 id, string memory tokenURI) external onlyOwner {
        _tokenURIs[id] = tokenURI;
        emit URI(uri(id), id);
    }

    function balanceOf(address account, uint256 id) public view override returns (uint256) {
        require(account != address(0), "ERC1155: balance query for the zero address");
        return _balances[id][account];
    }

    function balanceOfBatch(address[] memory accounts, uint256[] memory ids)
        public
        view
        override
        returns (uint256[] memory)
    {
        require(accounts.length == ids.length, "ERC1155: accounts and ids length mismatch");

        uint256[] memory batchBalances = new uint256[](accounts.length);

        for (uint256 i = 0; i < accounts.length; ++i) {
            batchBalances[i] = balanceOf(accounts[i], ids[i]);
        }

        return batchBalances;
    }

    function setApprovalForAll(address operator, bool approved) public override {
        require(msg.sender != operator, "ERC1155: setting approval status for self");
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function isApprovedForAll(address account, address operator) public view override returns (bool) {
        return _operatorApprovals[account][operator];
    }

    function safeTransferFrom(
        address from,
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) public override {
        require(
            from == msg.sender || isApprovedForAll(from, msg.sender),
            "ERC1155: caller is not owner nor approved"
        );
        _safeTransferFrom(from, to, id, amount, data);
    }

    function safeBatchTransferFrom(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) public override {
        require(
            from == msg.sender || isApprovedForAll(from, msg.sender),
            "ERC1155: transfer caller is not owner nor approved"
        );
        _safeBatchTransferFrom(from, to, ids, amounts, data);
    }

    function _safeTransferFrom(
        address from,
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) internal {
        require(to != address(0), "ERC1155: transfer to the zero address");

        uint256 fromBalance = _balances[id][from];
        require(fromBalance >= amount, "ERC1155: insufficient balance for transfer");

        unchecked {
            _balances[id][from] = fromBalance - amount;
        }
        _balances[id][to] += amount;

        emit TransferSingle(msg.sender, from, to, id, amount);

        _doSafeTransferAcceptanceCheck(msg.sender, from, to, id, amount, data);
    }

    function _safeBatchTransferFrom(
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal {
        require(ids.length == amounts.length, "ERC1155: ids and amounts length mismatch");
        require(to != address(0), "ERC1155: transfer to the zero address");

        for (uint256 i = 0; i < ids.length; ++i) {
            uint256 id = ids[i];
            uint256 amount = amounts[i];

            uint256 fromBalance = _balances[id][from];
            require(fromBalance >= amount, "ERC1155: insufficient balance for transfer");
            unchecked {
                _balances[id][from] = fromBalance - amount;
            }
            _balances[id][to] += amount;
        }

        emit TransferBatch(msg.sender, from, to, ids, amounts);

        _doSafeBatchTransferAcceptanceCheck(msg.sender, from, to, ids, amounts, data);
    }

    function _mint(
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) internal {
        require(to != address(0), "ERC1155: mint to the zero address");

        _balances[id][to] += amount;
        emit TransferSingle(msg.sender, address(0), to, id, amount);

        _doSafeTransferAcceptanceCheck(msg.sender, address(0), to, id, amount, data);
    }

    function _mintBatch(
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) internal {
        require(to != address(0), "ERC1155: mint to the zero address");
        require(ids.length == amounts.length, "ERC1155: ids and amounts length mismatch");

        for (uint256 i = 0; i < ids.length; i++) {
            _balances[ids[i]][to] += amounts[i];
        }

        emit TransferBatch(msg.sender, address(0), to, ids, amounts);

        _doSafeBatchTransferAcceptanceCheck(msg.sender, address(0), to, ids, amounts, data);
    }

    // 游戏相关功能
    function craftItem(uint256[] memory materialIds, uint256[] memory amounts, uint256 resultId) external {
        // 检查材料是否足够
        for (uint256 i = 0; i < materialIds.length; i++) {
            require(
                balanceOf(msg.sender, materialIds[i]) >= amounts[i],
                "Insufficient materials"
            );
        }

        // 消耗材料
        for (uint256 i = 0; i < materialIds.length; i++) {
            _balances[materialIds[i]][msg.sender] -= amounts[i];
        }

        // 给予制作结果
        _balances[resultId][msg.sender] += 1;

        emit TransferSingle(msg.sender, msg.sender, address(0), materialIds[0], amounts[0]);
        emit TransferSingle(msg.sender, address(0), msg.sender, resultId, 1);
    }

    function mint(address to, uint256 id, uint256 amount) external onlyOwner {
        _mint(to, id, amount, "");
    }

    function mintBatch(address to, uint256[] memory ids, uint256[] memory amounts) external onlyOwner {
        _mintBatch(to, ids, amounts, "");
    }

    function _doSafeTransferAcceptanceCheck(
        address operator,
        address from,
        address to,
        uint256 id,
        uint256 amount,
        bytes memory data
    ) private {
        if (to.code.length > 0) {
            try IERC1155Receiver(to).onERC1155Received(operator, from, id, amount, data) returns (bytes4 response) {
                if (response != IERC1155Receiver.onERC1155Received.selector) {
                    revert("ERC1155: ERC1155Receiver rejected tokens");
                }
            } catch Error(string memory reason) {
                revert(reason);
            } catch {
                revert("ERC1155: transfer to non ERC1155Receiver implementer");
            }
        }
    }

    function _doSafeBatchTransferAcceptanceCheck(
        address operator,
        address from,
        address to,
        uint256[] memory ids,
        uint256[] memory amounts,
        bytes memory data
    ) private {
        if (to.code.length > 0) {
            try IERC1155Receiver(to).onERC1155BatchReceived(operator, from, ids, amounts, data) returns (
                bytes4 response
            ) {
                if (response != IERC1155Receiver.onERC1155BatchReceived.selector) {
                    revert("ERC1155: ERC1155Receiver rejected tokens");
                }
            } catch Error(string memory reason) {
                revert(reason);
            } catch {
                revert("ERC1155: transfer to non ERC1155Receiver implementer");
            }
        }
    }

    // 支持的接口检查
    function supportsInterface(bytes4 interfaceId) public view virtual returns (bool) {
        return
            interfaceId == type(IERC1155).interfaceId ||
            interfaceId == type(IERC1155MetadataURI).interfaceId ||
            interfaceId == 0x01ffc9a7; // ERC165 interface ID
    }
}
```

## 8.4 简单 DeFi 协议

### 8.4.1 代币交换协议

```solidity
contract SimpleSwap {
    struct Pool {
        address tokenA;
        address tokenB;
        uint256 reserveA;
        uint256 reserveB;
        uint256 totalLiquidity;
        mapping(address => uint256) liquidity;
    }

    mapping(bytes32 => Pool) public pools;
    mapping(address => mapping(address => bytes32)) public getPoolId;

    event PoolCreated(address indexed tokenA, address indexed tokenB, bytes32 poolId);
    event LiquidityAdded(bytes32 indexed poolId, address indexed provider, uint256 amountA, uint256 amountB, uint256 liquidity);
    event LiquidityRemoved(bytes32 indexed poolId, address indexed provider, uint256 amountA, uint256 amountB, uint256 liquidity);
    event Swap(bytes32 indexed poolId, address indexed user, address tokenIn, uint256 amountIn, uint256 amountOut);

    function createPool(address tokenA, address tokenB) external returns (bytes32 poolId) {
        require(tokenA != tokenB, "Identical tokens");
        require(tokenA != address(0) && tokenB != address(0), "Zero address");

        if (tokenA > tokenB) {
            (tokenA, tokenB) = (tokenB, tokenA);
        }

        require(getPoolId[tokenA][tokenB] == 0, "Pool already exists");

        poolId = keccak256(abi.encodePacked(tokenA, tokenB));

        pools[poolId].tokenA = tokenA;
        pools[poolId].tokenB = tokenB;

        getPoolId[tokenA][tokenB] = poolId;
        getPoolId[tokenB][tokenA] = poolId;

        emit PoolCreated(tokenA, tokenB, poolId);
    }

    function addLiquidity(
        bytes32 poolId,
        uint256 amountADesired,
        uint256 amountBDesired,
        uint256 amountAMin,
        uint256 amountBMin
    ) external returns (uint256 amountA, uint256 amountB, uint256 liquidity) {
        Pool storage pool = pools[poolId];
        require(pool.tokenA != address(0), "Pool does not exist");

        if (pool.reserveA == 0 && pool.reserveB == 0) {
            // 首次添加流动性
            amountA = amountADesired;
            amountB = amountBDesired;
            liquidity = sqrt(amountA * amountB);
        } else {
            // 计算最优比例
            uint256 amountBOptimal = (amountADesired * pool.reserveB) / pool.reserveA;
            if (amountBOptimal <= amountBDesired) {
                require(amountBOptimal >= amountBMin, "Insufficient B amount");
                amountA = amountADesired;
                amountB = amountBOptimal;
            } else {
                uint256 amountAOptimal = (amountBDesired * pool.reserveA) / pool.reserveB;
                require(amountAOptimal <= amountADesired && amountAOptimal >= amountAMin, "Insufficient A amount");
                amountA = amountAOptimal;
                amountB = amountBDesired;
            }

            liquidity = min(
                (amountA * pool.totalLiquidity) / pool.reserveA,
                (amountB * pool.totalLiquidity) / pool.reserveB
            );
        }

        require(liquidity > 0, "Insufficient liquidity minted");

        // 转入代币
        IERC20(pool.tokenA).transferFrom(msg.sender, address(this), amountA);
        IERC20(pool.tokenB).transferFrom(msg.sender, address(this), amountB);

        // 更新储备和流动性
        pool.reserveA += amountA;
        pool.reserveB += amountB;
        pool.totalLiquidity += liquidity;
        pool.liquidity[msg.sender] += liquidity;

        emit LiquidityAdded(poolId, msg.sender, amountA, amountB, liquidity);
    }

    function removeLiquidity(
        bytes32 poolId,
        uint256 liquidity,
        uint256 amountAMin,
        uint256 amountBMin
    ) external returns (uint256 amountA, uint256 amountB) {
        Pool storage pool = pools[poolId];
        require(pool.liquidity[msg.sender] >= liquidity, "Insufficient liquidity");

        amountA = (liquidity * pool.reserveA) / pool.totalLiquidity;
        amountB = (liquidity * pool.reserveB) / pool.totalLiquidity;

        require(amountA >= amountAMin && amountB >= amountBMin, "Insufficient output amount");

        // 更新储备和流动性
        pool.reserveA -= amountA;
        pool.reserveB -= amountB;
        pool.totalLiquidity -= liquidity;
        pool.liquidity[msg.sender] -= liquidity;

        // 转出代币
        IERC20(pool.tokenA).transfer(msg.sender, amountA);
        IERC20(pool.tokenB).transfer(msg.sender, amountB);

        emit LiquidityRemoved(poolId, msg.sender, amountA, amountB, liquidity);
    }

    function swap(
        bytes32 poolId,
        address tokenIn,
        uint256 amountIn,
        uint256 amountOutMin
    ) external returns (uint256 amountOut) {
        Pool storage pool = pools[poolId];
        require(tokenIn == pool.tokenA || tokenIn == pool.tokenB, "Invalid token");

        bool isTokenA = tokenIn == pool.tokenA;
        (uint256 reserveIn, uint256 reserveOut) = isTokenA
            ? (pool.reserveA, pool.reserveB)
            : (pool.reserveB, pool.reserveA);

        require(amountIn > 0 && reserveIn > 0 && reserveOut > 0, "Invalid amounts");

        // 计算输出金额（扣除 0.3% 手续费）
        uint256 amountInWithFee = amountIn * 997;
        amountOut = (amountInWithFee * reserveOut) / (reserveIn * 1000 + amountInWithFee);

        require(amountOut >= amountOutMin, "Insufficient output amount");
        require(amountOut < reserveOut, "Insufficient liquidity");

        // 执行交换
        IERC20(tokenIn).transferFrom(msg.sender, address(this), amountIn);

        address tokenOut = isTokenA ? pool.tokenB : pool.tokenA;
        IERC20(tokenOut).transfer(msg.sender, amountOut);

        // 更新储备
        if (isTokenA) {
            pool.reserveA += amountIn;
            pool.reserveB -= amountOut;
        } else {
            pool.reserveB += amountIn;
            pool.reserveA -= amountOut;
        }

        emit Swap(poolId, msg.sender, tokenIn, amountIn, amountOut);
    }

    function getAmountOut(bytes32 poolId, address tokenIn, uint256 amountIn)
        external
        view
        returns (uint256 amountOut)
    {
        Pool storage pool = pools[poolId];
        require(tokenIn == pool.tokenA || tokenIn == pool.tokenB, "Invalid token");

        bool isTokenA = tokenIn == pool.tokenA;
        (uint256 reserveIn, uint256 reserveOut) = isTokenA
            ? (pool.reserveA, pool.reserveB)
            : (pool.reserveB, pool.reserveA);

        require(amountIn > 0 && reserveIn > 0 && reserveOut > 0, "Invalid amounts");

        uint256 amountInWithFee = amountIn * 997;
        amountOut = (amountInWithFee * reserveOut) / (reserveIn * 1000 + amountInWithFee);
    }

    function sqrt(uint256 x) internal pure returns (uint256) {
        if (x == 0) return 0;

        uint256 z = (x + 1) / 2;
        uint256 y = x;

        while (z < y) {
            y = z;
            z = (x / z + z) / 2;
        }

        return y;
    }

    function min(uint256 a, uint256 b) internal pure returns (uint256) {
        return a < b ? a : b;
    }
}
```

### 8.4.2 流动性挖矿

```solidity
contract LiquidityMining {
    IERC20 public rewardToken;
    IERC20 public stakingToken;

    uint256 public rewardRate = 100; // 每秒奖励代币数量
    uint256 public lastUpdateTime;
    uint256 public rewardPerTokenStored;

    mapping(address => uint256) public userRewardPerTokenPaid;
    mapping(address => uint256) public rewards;
    mapping(address => uint256) public balances;

    uint256 private _totalSupply;

    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount);
    event RewardPaid(address indexed user, uint256 reward);

    constructor(address _stakingToken, address _rewardToken) {
        stakingToken = IERC20(_stakingToken);
        rewardToken = IERC20(_rewardToken);
        lastUpdateTime = block.timestamp;
    }

    modifier updateReward(address account) {
        rewardPerTokenStored = rewardPerToken();
        lastUpdateTime = block.timestamp;

        if (account != address(0)) {
            rewards[account] = earned(account);
            userRewardPerTokenPaid[account] = rewardPerTokenStored;
        }
        _;
    }

    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) external view returns (uint256) {
        return balances[account];
    }

    function rewardPerToken() public view returns (uint256) {
        if (_totalSupply == 0) {
            return rewardPerTokenStored;
        }

        return rewardPerTokenStored +
            ((block.timestamp - lastUpdateTime) * rewardRate * 1e18) / _totalSupply;
    }

    function earned(address account) public view returns (uint256) {
        return (balances[account] * (rewardPerToken() - userRewardPerTokenPaid[account])) / 1e18 + rewards[account];
    }

    function stake(uint256 amount) external updateReward(msg.sender) {
        require(amount > 0, "Cannot stake 0");

        _totalSupply += amount;
        balances[msg.sender] += amount;

        stakingToken.transferFrom(msg.sender, address(this), amount);

        emit Staked(msg.sender, amount);
    }

    function withdraw(uint256 amount) public updateReward(msg.sender) {
        require(amount > 0, "Cannot withdraw 0");
        require(balances[msg.sender] >= amount, "Insufficient balance");

        _totalSupply -= amount;
        balances[msg.sender] -= amount;

        stakingToken.transfer(msg.sender, amount);

        emit Withdrawn(msg.sender, amount);
    }

    function getReward() public updateReward(msg.sender) {
        uint256 reward = rewards[msg.sender];
        if (reward > 0) {
            rewards[msg.sender] = 0;
            rewardToken.transfer(msg.sender, reward);
            emit RewardPaid(msg.sender, reward);
        }
    }

    function exit() external {
        withdraw(balances[msg.sender]);
        getReward();
    }

    // 管理员函数
    function setRewardRate(uint256 _rewardRate) external updateReward(address(0)) {
        rewardRate = _rewardRate;
    }
}
```

## 8.5 代币安全最佳实践

### 8.5.1 安全检查和限制

```solidity
contract SecureToken is AdvancedERC20 {
    uint256 public maxTransferAmount;
    uint256 public dailyTransferLimit;
    mapping(address => uint256) public lastTransferDay;
    mapping(address => uint256) public dailyTransferred;

    // 防止三明治攻击的时间锁
    mapping(address => uint256) public lastApprovalTime;
    uint256 public constant APPROVAL_DELAY = 1 hours;

    // 防重入锁
    bool private locked;

    modifier noReentrant() {
        require(!locked, "ReentrancyGuard: reentrant call");
        locked = true;
        _;
        locked = false;
    }

    modifier validTransferAmount(uint256 amount) {
        require(amount <= maxTransferAmount, "Transfer amount exceeds limit");
        _;
    }

    modifier checkDailyLimit(address from, uint256 amount) {
        uint256 currentDay = block.timestamp / 1 days;

        if (lastTransferDay[from] < currentDay) {
            dailyTransferred[from] = 0;
            lastTransferDay[from] = currentDay;
        }

        require(
            dailyTransferred[from] + amount <= dailyTransferLimit,
            "Daily transfer limit exceeded"
        );

        dailyTransferred[from] += amount;
        _;
    }

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        uint256 totalSupply_,
        uint256 _maxTransferAmount,
        uint256 _dailyTransferLimit
    ) AdvancedERC20(name_, symbol_, decimals_, totalSupply_) {
        maxTransferAmount = _maxTransferAmount;
        dailyTransferLimit = _dailyTransferLimit;
    }

    function _transfer(address from, address to, uint256 amount)
        internal
        override
        noReentrant
        validTransferAmount(amount)
        checkDailyLimit(from, amount)
    {
        super._transfer(from, to, amount);
    }

    function approve(address spender, uint256 amount) public override returns (bool) {
        // 防止前置攻击
        require(
            allowance(msg.sender, spender) == 0 || amount == 0,
            "Must reset allowance to 0 first"
        );

        lastApprovalTime[msg.sender] = block.timestamp;
        return super.approve(spender, amount);
    }

    function transferFrom(address from, address to, uint256 amount) public override returns (bool) {
        // 检查时间锁
        require(
            block.timestamp >= lastApprovalTime[from] + APPROVAL_DELAY,
            "Approval time lock not expired"
        );

        return super.transferFrom(from, to, amount);
    }

    // 紧急暂停功能
    function emergencyPause() external onlyOwner {
        setPaused(true);
    }

    // 批量转账（谨慎使用）
    function batchTransfer(address[] calldata recipients, uint256[] calldata amounts)
        external
        noReentrant
        returns (bool)
    {
        require(recipients.length == amounts.length, "Arrays length mismatch");
        require(recipients.length <= 200, "Too many recipients");

        uint256 totalAmount = 0;
        for (uint256 i = 0; i < amounts.length; i++) {
            totalAmount += amounts[i];
        }

        require(balanceOf(msg.sender) >= totalAmount, "Insufficient balance");

        for (uint256 i = 0; i < recipients.length; i++) {
            _transfer(msg.sender, recipients[i], amounts[i]);
        }

        return true;
    }

    // 设置转账限额
    function setMaxTransferAmount(uint256 _maxTransferAmount) external onlyOwner {
        maxTransferAmount = _maxTransferAmount;
    }

    function setDailyTransferLimit(uint256 _dailyTransferLimit) external onlyOwner {
        dailyTransferLimit = _dailyTransferLimit;
    }
}
```

## 8.6 实战练习

### 练习 1：创建自己的代币

创建一个具有以下特性的 ERC-20 代币：

- 固定总供应量
- 可暂停转账
- 黑名单功能
- 代币销毁

### 练习 2：NFT 收藏品

设计一个 NFT 收藏品合约：

- 限量发行
- 白名单预售
- 公开销售
- 版税功能

### 练习 3：简单 DEX

实现一个简单的去中心化交易所：

- 代币交换
- 流动性提供
- 手续费分配
- 价格预言机

## 8.7 章节总结

在本章中，我们深入学习了：

1. **ERC-20 标准**：同质化代币的完整实现
2. **ERC-721 标准**：NFT 的创建和管理
3. **ERC-1155 标准**：多代币标准的应用
4. **DeFi 基础**：代币交换和流动性挖矿
5. **安全实践**：代币合约的安全考虑

### 关键知识点

- 代币标准定义了通用的接口规范
- NFT 提供了数字资产的唯一性证明
- DeFi 协议基于代币实现金融服务
- 安全性是代币合约的首要考虑

### 下一章预告

在下一章《安全性与最佳实践》中，我们将深入学习：

- 常见的智能合约安全漏洞
- 防范攻击的技术手段
- 安全审计的方法和工具
- 代码质量保证策略

继续您的 Solidity 学习之旅：[第九章 - 安全性与最佳实践](09_security_best_practices.md)
