# Project ARIADNE — Prototype Phase 5：Windows Application (v2 Draft)

> **Status: DRAFT**  
> 本文書はPhase 5正本再構成（Issue #22）の作業用ドラフトである。A02〜A06で段階的に整備し、A07で正式切替が承認されるまで、現行正本 `phase-5-windows-application.md` を置き換えない。

## 1. Purpose & Scope

### 1.1 Purpose

Prototype Phase 5では、Svelte / Wails v3 / GoによるWindows Applicationを実装し、PROJECT ARIADNEのWindows Applicationとしての技術的成立性を検証する。

Task管理Applicationを題材に、Frontend、Wails Boundary、Application Service、File Persistenceを連携させ、Taskの参照・作成・編集・永続化・復元までを実際に動作させる。

本Phaseの目的は、Applicationが動作することの確認だけではない。採用したApplication Architectureについて、責務分離と境界が設計どおりに実装可能であることに加え、**その設計判断自体が妥当であるか**を実装・検証を通じて確かめる。

### 1.2 Design Principles

Phase 5では、以下の4原則に基づいて設計・開発する。

1. **理解と納得に基づく設計**  
   採用する設計の理由と選択肢を理解し、納得したうえで決定する。テルとアリアドネは、提案・疑問・反論を交わしながら設計判断を深める。既存の実装が動くことだけを理由に判断を固定しない。
2. **全体の整合性を守る**  
   個々の機能や局所的な実装の都合だけでなく、Architecture全体の責務境界、採用済み設計判断、関連文書との整合性を重視する。
3. **小さく成立させ、検証を通じて育てる**  
   完成形を一度に作り込まず、小さな単位で設計・実装・Test・Build・実機確認を繰り返す。未知の技術要素は段階的に検証し、その結果を次の判断に反映する。
4. **美しく、無理のない構造を選ぶ**  
   不必要な複雑さや先回りした抽象化を避け、責務が明確で、自然に成立する構造を選ぶ。実装コストだけでなく、構造のわかりやすさと将来の変更容易性を考慮する。

具体的な開発サイクル、Test・Buildの運用、設計判断の記録方法および協働の進め方は、第2章 Development Processで扱う。

### 1.3 Scope

Prototype Phase 5では、Task管理を題材として、以下を対象とする。

- Taskの一覧・詳細参照、新規作成、編集
- ローカルYAMLへのTask保存、およびApplication再起動後の復元
- TaskのValidationと、正常系・異常系のApplication内での伝達・処理
- Windows ApplicationとしてのUI、Build、Test、実機での起動・操作確認
- Frontend / Wails Boundary / Application Service / File Persistenceの連携、責務分離およびArchitecture設計判断の検証

Taskの**新規作成**は対象に含む。一方、Task管理環境の**初回初期化**（`tasks.yaml` が存在しない場合の初期化を含む）は対象外とし、Coreで扱う。Taskが0件であること自体は正常な状態として扱う。

Phase 2〜4で扱ったElements / API / DDLのモデルおよび変換構造を、Applicationから操作する機能は本Phaseの対象外とする。それらを扱うApplication機能やValidator / Resolver / Generator等の実装はCoreで扱う。

### 1.4 Completion Criteria

Phase 5の完了は、以下の観点で判断する。

1. **Applicationの成立**：Task管理の対象操作とYAML永続化・復元が、Windows Applicationとして動作する。
2. **Architectureの成立**：採用した責務分離と境界に基づいてApplicationを構成でき、その設計判断の妥当性を検証できている。
3. **検証の成立**：必要な自動Test、BuildおよびWindows実機での動作確認が成立している。
4. **判断の整理**：検証を通じて確定した重要な設計判断、残課題およびCoreへの引継ぎ事項が整理されている。

個別BehaviorやTestの詳細な合否条件は、第2〜4章および関連するBehavior Map・Test文書で定義する。
